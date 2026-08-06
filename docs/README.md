# Scind User Documentation — placeholder

This directory is reserved for Scind's **user-facing documentation track**. It is
empty today. No user-facing content exists yet.

## Intent

Scind keeps two documentation tracks, one per audience:

| Track | Directory | Audience | Organization |
|-------|-----------|----------|--------------|
| Engineering canon | [`engineering/`](../engineering/README.md) | contributors, agents, spec readers | Layered Documentation System (LDS): decisions, product, architecture, specs, reference, implementation |
| User documentation | `docs/` (this directory) | people who use Scind | [Diátaxis](https://diataxis.fr/): tutorials, how-to guides, reference, explanation |

The engineering track answers "why is Scind built this way, and what exactly must
it do". The user track answers "how do I get my work done with Scind". Mixing the
two makes both harder to read, so they stay apart.

## Status

- `engineering/` — populated. It held the path `docs/` until the rename recorded
  below.
- `docs/` — this placeholder only.

Write the Diátaxis content here when Scind has users to write it for. Until then,
send readers to [`engineering/`](../engineering/README.md).

## History

The engineering canon moved from `docs/` to `engineering/` so that the directory
name states its audience, and so that Scind and Xcind expose the same comparison
surface (`scind/engineering/` ↔ `xcind/engineering/`) for cross-project sync.

Sync artifacts written before the rename refer to Scind paths as `docs/…`. Read
those as `engineering/…`. See the path-alias rule in Xcind's
`engineering/maintenance/cross-project-sync.md`.

**Tracking**: reconciliation-ledger row `RL-038` / cross-project finding `XA-0043`
(two-track documentation).
