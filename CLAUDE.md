This projects exists soley to work through the 
This directory exists for the sole purpose of working through issues identified during a review of the Contrail specification documents.

When working on anything, always read the following for proper context:

- @README.md – contains an overview of the purpose of this document
- @specs – the entire collection of Contrail specification documents

Then, read @issues/00-index.md to get an overview of the issues.

When asked to work on the next issue, refer to the list of completed issues in the issue index file (if it exists), and choose the first issue from the recommended list that is not already completed.

When reading an issue, if any of the responses are missing or are still set to "_[Your response here]_", show the details of the question and ask for my response before working on that specific task. Then, update the issue with my response as a block quote and replace the original "_[Your response here]_" text.

## Version Management

When modifying spec documents in the `specs/` directory:

1. **Increment versions independently** — Each document maintains its own version number
2. **Update at session end** — At the end of a session where changes are made to spec documents, increment the patch version (e.g., 0.5.0 → 0.5.1) for each modified document
3. **Update revision history** — Add an entry to the Revision History table at the bottom of each modified document with a brief description of changes

Example revision history entry:
```markdown
| 0.5.1-draft | Dec 2024 | Removed logs command, deferred to contrail-compose |
```