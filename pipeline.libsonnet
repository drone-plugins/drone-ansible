{
  build(name='', os='linux', arch='amd64', version='')::
    local tag = if os == 'windows' then os + '-' + version else os + '-' + arch;
    local file_suffix = std.strReplace(tag, '-', '.');
    local volumes = if os == 'windows' then [{ name: 'docker_pipe', path: '//./pipe/docker_engine' }] else [];
    local gopath = { name: 'gopath', temp: {} };

    {
      kind: 'pipeline',
      name: tag,
      platform: {
        os: os,
        arch: arch,
        version: if std.length(version) > 0 then version,
      },
      steps: [
        {
          name: 'build-push',
          image: 'golang:1.11',
          pull: 'always',
          environment: {
            CGO_ENABLED: '0',
            GO111MODULE: 'on',
          },
          commands: [
            'go build -v -ldflags "-X main.build=${DRONE_BUILD_NUMBER}" -a -o release/' + os + '/' + arch + '/drone-' + name,
          ],
          volumes: [
            {
              name: 'gopath',
              path: '/go',
            }
          ],
          when: {
            event: ['push', 'pull_request'],
          },
        },
        {
          name: 'build-tag',
          image: 'golang:1.11',
          pull: 'always',
          environment: {
            CGO_ENABLED: '0',
            GO111MODULE: 'on',
          },
          commands: [
            'go build -v -ldflags "-X main.version=${DRONE_TAG##v} -X main.build=${DRONE_BUILD_NUMBER}" -a -o release/' + os + '/' + arch + '/drone-' + name,
          ],
          volumes: [
            {
              name: 'gopath',
              path: '/go',
            }
          ],
          when: {
            event: ['tag'],
          },
        },
        {
          name: 'executable',
          image: 'golang:1.11',
          pull: 'always',
          commands: [
            './release/' + os + '/' + arch + '/drone-' + name + ' --help',
          ],
        },
        {
          name: 'dryrun',
          image: 'plugins/docker:' + tag,
          pull: 'always',
          settings: {
            dry_run: true,
            tags: tag,
            dockerfile: 'docker/Dockerfile.' + file_suffix,
            repo: 'plugins/' + name,
            username: {
              from_secret: 'docker_username',
            },
            password: {
              from_secret: 'docker_password',
            },
          },
          volumes: if std.length(volumes) > 0 then volumes,
          when: {
            event: ['pull_request'],
          },
        },
        {
          name: 'publish',
          image: 'plugins/docker:' + tag,
          pull: 'always',
          settings: {
            auto_tag: true,
            auto_tag_suffix: tag,
            dockerfile: 'docker/Dockerfile.' + file_suffix,
            repo: 'plugins/' + name,
            username: {
              from_secret: 'docker_username',
            },
            password: {
              from_secret: 'docker_password',
            },
          },
          volumes: if std.length(volumes) > 0 then volumes,
          when: {
            event: ['push', 'tag'],
          },
        },
      ],
      depends_on: [
        'testing',
      ],
      trigger: {
        branch: ['master'],
      },
      volumes: if os == 'windows' then [gopath, { name: 'docker_pipe', host: { path: '//./pipe/docker_engine' } }] else [gopath],
    },

  notifications(os='linux', arch='amd64', version='', depends_on=[])::
    {
      kind: 'pipeline',
      name: 'notifications',
      platform: {
        os: os,
        arch: arch,
        version: if std.length(version) > 0 then version,
      },
      steps: [
        {
          name: 'manifest',
          image: 'plugins/manifest:1',
          pull: 'always',
          settings: {
            auto_tag: true,
            username: {
              from_secret: 'docker_username',
            },
            password: {
              from_secret: 'docker_password',
            },
            spec: 'docker/manifest.tmpl',
            ignore_missing: true,
          },
        },
        {
          name: 'microbadger',
          image: 'plugins/webhook:1',
          pull: 'always',
          settings: {
            url: {
              from_secret: 'microbadger_url',
            },
          },
        },
      ],
      depends_on: depends_on,
      trigger: {
        event: ['push', 'tag'],
      },
    },
}
