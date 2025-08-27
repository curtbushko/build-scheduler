# build scheduler

Basis for idea: [Build Queues on Vercel](https://vercel.com/docs/builds/build-queues)

From the doc:

# Concurrency queue
This queue manages how many builds can run in parallel based on the number of concurrent build slots available to the team and user. If all concurrent build slots are in use, new builds are queued until a slot becomes available unless you have On-Demand Concurrent Builds enabled at the project level.

## How concurrent build slots work

Concurrent build slots are the key factor in concurrent build queuing. They control how many builds can run at the same time and ensure efficient use of resources while prioritizing the latest changes.

Each account plan comes with a predefined number of build slots:

- Hobby accounts allow one build at a time.
- Pro accounts support up to 12 simultaneous builds.
- Enterprise accounts can have custom limits based on their plan.

# Git branch queue

When multiple builds are pushed to the same Git branch, they are handled sequentially. If new commits are pushed while a build is in progress:

1. The current build is completed first.
2. Queued builds for earlier commits are skipped.
3. The most recent commit is built and deployed.

This ensures that only the latest changes are deployed, reducing unnecessary builds and improving deployment efficiency.

# Concurrencty on the same branch

It's possible for builds to be affected by both queues simultaneously. For example, if your builds are queued due to no slots being available and you submitted multiple commits to the same branch, the following will happen:

- The latest commit to that branch will start building when a slot becomes available.
- All previous commits to that branch will be skipped when that happens.


# TODO
- [ ] Create a build type
    - should have: buildID, teamID, buildStart, buildEnd, queueStart, queueEnd, repo, branch, commit, status
- [ ] Create a Team type
    - should have a []builds queue
    - teamID, totalQueueTime, totalBuildTime, plan, builds_active

Worker Nodes can be self regulating and registering!
