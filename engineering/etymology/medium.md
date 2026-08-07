# The Etymology of Scind

## The Latin Root

**Scind** derives from the Latin verb *scindere*, meaning "to cut," "to split," or "to tear apart." The word belongs to a family of terms that describe decisive separation — not gradual erosion or blending, but a clean division into distinct parts.

From *scindere*, Latin gave us *scissor* (one who cuts), *rescindere* (to cut back, to annul), and *abscindere* (to cut off). These passed into English as *scissors*, *rescind*, and *abscission*. The Greek cognate *skhizein* (to split) produced *schism* and *schizophrenia*. Every descendant carries the same core meaning: something whole becomes cleanly separated.

## The Connection to the Tool

Scind is a workspace orchestration system for Docker Compose. Its fundamental operation is taking a multi-application stack — frontend, backend, database, and other services that work together — and splitting it into multiple isolated instances that run simultaneously on a single host.

This maps directly to what *scindere* describes:

**Workspace isolation.** Each workspace is a clean cut. A `dev` workspace and a `review` workspace contain the same applications but share no state — separate containers, separate volumes, separate networks. One can be torn down without affecting the other.

**Pure overlay design.** Scind separates integration from application. Your Docker Compose files remain untouched; all workspace-specific configuration lives in generated override files layered on top. The application and the workspace concerns are split apart, each maintaining its own integrity.

**Two-layer networking.** External traffic flows through one network (the proxy layer); internal communication flows through another (the workspace-internal network). This is a structural split: the boundary between "how the outside world reaches your services" and "how your services talk to each other" is explicit and enforced.

**Convention-based naming.** Every hostname, alias, and project name follows a deterministic pattern. There is no ambiguity about which workspace owns which container. The naming conventions act as clean incision lines — they make the separations visible and predictable.

## Why This Name

The name *Scind* was chosen for the same reasons the tool was built: precision and economy.

It is five letters and one syllable. It is a direct imperative — a command, not a description. It does not explain what it does; it performs the action. In this way it aligns with the Unix tradition of short, sharp tool names that say just enough: `cut`, `split`, `grep`, `sed`.

The Latin root also carries a connotation that matters: *scindere* implies a deliberate act, not an accident. Scind does not merely happen to isolate workspaces. It is designed, from its conventions to its generated configuration, to make separation the default and entanglement impossible.
