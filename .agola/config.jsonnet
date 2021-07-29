local go_runtime(version, arch) = {
  type: 'pod',
  arch: arch,
  containers: [
    { image: 'golang:' + version },
  ],
};

local task_build_go(setup) = {
  name: 'build go ' + setup.goos,
  runtime: go_runtime('1.3', 'amd64'),
  working_dir: '/go/src/github.com/ercole-io/ercole-agent-rhel5',
  environment: {
    GOOS: setup.goos,
    BIN: setup.bin,
  },
  steps: [
    { type: 'clone' },
    {
      type: 'run',
      name: 'build',
      command: |||
        if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
          export VERSION=latest
          export BUILD_VERSION=${AGOLA_GIT_COMMITSHA}
        else
          export VERSION=${AGOLA_GIT_TAG}
          export BUILD_VERSION=${AGOLA_GIT_TAG}
        fi

        echo VERSION: ${VERSION}
        echo BUILD_VERSION: ${BUILD_VERSION}

        go build -ldflags="-X main.version=${BUILD_VERSION}" -o ${BIN}
      |||,
    },
    {
      type: 'save_to_workspace',
      contents: [{
        source_dir: '.',
        dest_dir: '.',
        paths: [
          setup.bin,
          'Makefile',
          'package/**',
          'fetch/**',
          'sql/**',
          'config.json',
          'LICENSE',  // Needed by windows
        ],
      }],
    },
  ],
  depends: ['test'],
};

local task_pkg_build_rhel(setup) = {
  name: 'pkg build ' + setup.dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: setup.pkg_build_image },
    ],
  },
  working_dir: '/project',
  environment: {
    WORKSPACE: '/project',
    DIST: setup.dist,
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'version',
      command: |||
        if [ -z ${AGOLA_GIT_TAG} ] || [[ ${AGOLA_GIT_TAG} == *-* ]]; then 
          export VERSION=latest
        else
          export VERSION=${AGOLA_GIT_TAG}
        fi
        echo VERSION: ${VERSION}
        echo "export VERSION=${VERSION}" > /tmp/variables
      |||,
    },
    {
      type: 'run',
      name: 'sed version',
      command: |||
        source /tmp/variables

        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel8/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel7/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel6/ercole-agent.spec
        sed -i "s|ERCOLE_VERSION|${VERSION}|g" package/rhel5/ercole-agent.spec
      |||,
    },
    { type: 'run', command: 'if [ $DIST == "rhel5" ]; then echo \'%_topdir %(echo \\$HOME)/rpmbuild\' > ~/.rpmmacros ; fi' },
    { type: 'run', command: 'rpmbuild --quiet -bl package/${DIST}/ercole-agent.spec || echo ok' },
    { type: 'run', command: 'if [ $DIST == "rhel5" ]; then mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS} ; fi' },
    { type: 'run', command: 'source /tmp/variables && mkdir -p ~/rpmbuild/SOURCES/ercole-agent-${VERSION}' },
    { type: 'run', command: 'source /tmp/variables && cp -r * ~/rpmbuild/SOURCES/ercole-agent-${VERSION}/' },
    { type: 'run', command: 'source /tmp/variables && tar -C ~/rpmbuild/SOURCES -cvzf ~/rpmbuild/SOURCES/ercole-agent-${VERSION}.tar.gz ercole-agent-${VERSION}' },
    { type: 'run', command: 'pwd; rpmbuild -v -bb package/${DIST}/ercole-agent.spec' },
    { type: 'run', command: 'mkdir dist' },
    { type: 'run', command: 'source /tmp/variables && cd ${WORKSPACE} && cp ~/rpmbuild/RPMS/x86_64/ercole-agent-${VERSION}-1*.x86_64.rpm dist/' },
    { type: 'run', command: 'cp ~/rpmbuild/RPMS/x86_64/ercole-*.rpm ${WORKSPACE}/dist' },
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: ['build go linux'],
};

local task_deploy_repository(dist) = {
  name: 'deploy repository.ercole.io ' + dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: 'curlimages/curl' },
    ],
  },
  environment: {
    REPO_USER: { from_variable: 'repo-user' },
    REPO_TOKEN: { from_variable: 'repo-token' },
    REPO_UPLOAD_URL: { from_variable: 'repo-upload-url' },
    REPO_INSTALL_URL: { from_variable: 'repo-install-url' },
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'curl',
      command: |||
        cd dist

        for f in *; do
          mv $f ${f/x86_64/el5.x86_64}
        done

        for f in *; do
        	URL=$(curl --user "${REPO_USER}" \
            --upload-file $f ${REPO_UPLOAD_URL} --insecure)
        	echo $URL
        	md5sum $f
        	curl -H "X-API-Token: ${REPO_TOKEN}" \
          -H "Content-Type: application/json" --request POST --data "{ \"filename\": \"$f\", \"url\": \"$URL\" }" \
          ${REPO_INSTALL_URL} --insecure
        done
      |||,
    },
  ],
  depends: ['pkg build ' + dist],
  when: {
    tag: '#.*#',
    branch: 'master',
  },
};

{
  runs: [
    {
      name: 'ercole-agent-rhel5',
      tasks: [
        { name: 'test',
          runtime: go_runtime('1.3', 'amd64'),
          working_dir: '/go/src/github.com/ercole-io/ercole-agent-rhel5',
          steps: [
            { type: 'clone' },
            { type: 'run', name: '', command: 'go test ./...' },
          ],
        },
      ] + [
        task_build_go(setup)
        for setup in [
          { goos: 'linux', bin: 'ercole-agent' },
        ]
      ] + [
        task_pkg_build_rhel(setup)
        for setup in [
          { pkg_build_image: 'amreo/rpmbuild-centos5', dist: 'rhel5', distfamily: 'rhel' },
        ]
      ] + [
        task_deploy_repository(dist)
        for dist in ['rhel5']
      ],
    },
  ],
}
