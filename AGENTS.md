# Copilot Instructions — Thoth (Daltu · Scribe · Stylus)

## Golden Rules

- Always prioritize security, privacy, and data integrity.
- Use clear, concise commit messages following Conventional Commits.
- Validate all code changes with tests before committing.
- Maintain consistent style and formatting across Daltu, Scribe, and Stylus.
- Preserve backward compatibility unless explicitly instructed otherwise.

---

## Scope

These instructions guide AI agents in automating development and maintenance tasks within the Thoth ecosystem, encompassing:

- Daltu: Core computation and logic.
- Scribe: Documentation and narrative generation.
- Stylus: Theming and styling components.

Agents must collaborate seamlessly, respecting module boundaries and shared conventions.

---

## Task Workflow

1. **Planning:** Break down high-level requests into subtasks with clear objectives.
2. **Implementation:** Generate minimal, test-covered code changes aligned with subtasks.
3. **Review:** Critically analyze code for correctness, security, and style compliance.
4. **Memory Update:** Record new rules or preferences in the Memory section for future reference.

---

## Branch/Commit Conventions

- Branches:
  - Feature branches: `feat/<short-description>`
  - Bugfix branches: `bug/<short-description>`
  - Documentation: `docs/<topic>`
- Commits:
  - Use Conventional Commit format, e.g., `feat: add R-value policy enforcement`
  - Include references to tasks or memory entries when applicable.
  - Ensure commit messages are imperative and descriptive.

---

## Memory Function

- Agents must document new rules, preferences, or decisions in the Memory section of AGENTS.md.
- Memory entries follow the format `[MEM-XXXX] YYYY-MM-DD — Summary`.
- Retired or obsolete memories should be moved to the Retired Memories subsection with proper annotation.
- All memory updates must be committed with the format `memory(<topic>): <summary>`.
- Memory updates should only apply to instruction-related files:
  - `AGENTS.md`
  - `CLINE.md`
  - `.github/copilot-instructions.md`

---

## Architecture & Security Rules

- Enforce modular architecture boundaries between Daltu, Scribe, and Stylus.
- Validate all inputs and sanitize outputs to prevent injection or leakage.
- Use cryptographic hashing and durable storage patterns for sensitive data.
- Avoid introducing breaking changes without explicit approval.
- Maintain audit trails for all automated changes.

---

## Testing

- Unit tests must cover all new logic paths.
- Integration tests should verify agent handoffs and data flow.
- Behavioral tests ensure Memory updates persist and apply correctly.
- End-to-end tests confirm agents produce merge-ready pull requests.
- Agents must regenerate code if tests fail before committing.

---

## Documentation

- Update relevant markdown files to reflect feature additions or changes.
- Ensure examples and usage instructions are clear and accurate.
- Link to related tasks and memory entries for traceability.
- Use consistent terminology and formatting throughout.

---

## Example Policies

- Always reject code that weakens durability or privacy guarantees.
- Enforce commit and branch naming conventions strictly.
- Require tests accompany any functional code changes.
- Document any exceptions or overrides in the Memory section.
- Log all agent actions in structured JSON format under `/logs`.

---

## End of Instructions