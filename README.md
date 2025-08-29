# BUILD SCHEDULER

Basis for idea: [Build Queues on Vercel](https://vercel.com/docs/builds/build-queues)

From the website:

## Concurrency queue
This queue manages how many builds can run in parallel based on the number of concurrent build slots available to the team and user. If all concurrent build slots are in use, new builds are queued until a slot becomes available unless you have On-Demand Concurrent Builds enabled at the project level.

## How concurrent build slots work

Concurrent build slots are the key factor in concurrent build queuing. They control how many builds can run at the same time and ensure efficient use of resources while prioritizing the latest changes.

Each account plan comes with a predefined number of build slots:

- Hobby accounts allow one build at a time.
- Pro accounts support up to 12 simultaneous builds.
- Enterprise accounts can have custom limits based on their plan.

## Git branch queue

When multiple builds are pushed to the same Git branch, they are handled sequentially. If new commits are pushed while a build is in progress:

1. The current build is completed first.
2. Queued builds for earlier commits are skipped.
3. The most recent commit is built and deployed.

This ensures that only the latest changes are deployed, reducing unnecessary builds and improving deployment efficiency.

## Concurrencty on the same branch

It's possible for builds to be affected by both queues simultaneously. For example, if your builds are queued due to no slots being available and you submitted multiple commits to the same branch, the following will happen:

- The latest commit to that branch will start building when a slot becomes available.
- All previous commits to that branch will be skipped when that happens.

# BUILDING AND RUNNING

The environment is setup using nix. Everything uses make targets.

- 'make build' will build a 'build-scheduler' binary bin directory
- 'make run' will run the 'build-scheduler' from the bin directory

# EXAMPLE OUTPUT

```
Starting dispatcher
Worker started 3
Worker started 2
Worker started 4
Submitting build to team queue 58a7fc09-3fda-4e15-ae12-11080e3c40ba
Submitting build to team queue ccfe6657-6717-4bc0-afda-5ad7149c0090
Submitting build to team queue 711458b4-59b2-40bb-ac5c-3e9b25ad2697
Submitting build to team queue 97155754-5f28-4a92-98b0-9277608577dc
Submitting build to team queue d00b5b0f-3bd5-4fc7-9794-766e2a86d1bc
Worker started 1
Running build 58a7fc09-3fda-4e15-ae12-11080e3c40ba
Running build ccfe6657-6717-4bc0-afda-5ad7149c0090
Build complete 58a7fc09-3fda-4e15-ae12-11080e3c40ba
Build completed: ID=58a7fc09-3fda-4e15-ae12-11080e3c40ba, Status=complete, Team=1234, Duration=2562047h47m16.854775807s
Running build 711458b4-59b2-40bb-ac5c-3e9b25ad2697
Build complete ccfe6657-6717-4bc0-afda-5ad7149c0090
Build completed: ID=ccfe6657-6717-4bc0-afda-5ad7149c0090, Status=complete, Team=1234, Duration=2562047h47m16.854775807s
Running build 97155754-5f28-4a92-98b0-9277608577dc
Build complete 711458b4-59b2-40bb-ac5c-3e9b25ad2697
Build completed: ID=711458b4-59b2-40bb-ac5c-3e9b25ad2697, Status=complete, Team=1234, Duration=2562047h47m16.854775807s
Running build d00b5b0f-3bd5-4fc7-9794-766e2a86d1bc
Build complete 97155754-5f28-4a92-98b0-9277608577dc
Build completed: ID=97155754-5f28-4a92-98b0-9277608577dc, Status=complete, Team=1234, Duration=2562047h47m16.854775807s
Build complete d00b5b0f-3bd5-4fc7-9794-766e2a86d1bc
Build completed: ID=d00b5b0f-3bd5-4fc7-9794-766e2a86d1bc, Status=complete, Team=1234, Duration=2562047h47m16.854775807s

```

# TODO
- [x] Create a build type
    - should have: buildID, teamID, buildStart, buildEnd, queueStart, queueEnd, repo, branch, commit, status
- [x] Create a Team type
    - should have a []builds queue
    - teamID, totalQueueTime, totalBuildTime, plan, builds_active
- [x] Create dispatcher
- [x] Dispatch builds to workers based on slots
- [ ] Add tests

Worker Nodes can be self regulating and registering!
