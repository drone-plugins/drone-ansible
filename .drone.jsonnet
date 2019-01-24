local pipeline = import 'pipeline.libsonnet';

local PipelineTesting = {
  kind: 'pipeline',
  name: 'testing',
  platform: {
    os: 'linux',
    arch: 'amd64',
  },
  steps: [
    {
      name: 'vet',
      image: 'golang:1.11',
      pull: 'always',
      environment: {
        GO111MODULE: 'on',
      },
      commands: [
        'go vet ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        }
      ],
    },
    {
      name: 'test',
      image: 'golang:1.11',
      pull: 'always',
      environment: {
        GO111MODULE: 'on',
      },
      commands: [
        'go test -cover ./...',
      ],
      volumes: [
        {
          name: 'gopath',
          path: '/go',
        }
      ],
    },
  ],
  volumes: [
    {
      name: 'gopath',
      temp: {},
    },
  ],
};

[
  PipelineTesting,
  pipeline.build('ansible', 'linux', 'amd64'),
  pipeline.build('ansible', 'linux', 'arm64'),
  pipeline.build('ansible', 'linux', 'arm'),
  pipeline.notifications(depends_on=[
    'linux-amd64',
    'linux-arm64',
    'linux-arm',
  ]),
]
