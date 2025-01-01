# chfs

concurrent hash-trie filesystem

it's sort of like btrfs, but with a few small additions and
subtractions for my sake.

A filesystem provided via a git-like and fuse interface that is

* **versioned** it tracks a history of your filesystem on a
  per-write-event basis
* **lazy** unlike git, you don't need to download all files to get a
  specific version
* **verifiable** you can independently validate the integrity of any
  subset of the filesystem

## Why?

While version control has proven to be fairly useful for tracking
code, tracking it tends to get a bit more difficult once you get to
tracking massive amounts of data. This leads to folks designing
version control inside their database systems, which creates
unnecessary heft.

## Usage

### Integrity Verification
