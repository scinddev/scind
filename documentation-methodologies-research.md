# Comprehensive Documentation Methodologies Research

**Research Date**: January 2026
**Purpose**: Identify methodologies and standards for design-first, living, and reference documentation for developer-focused projects
**Application**: See `DOCUMENTATION-METHODOLOGY-ANALYSIS.md` for Contrail-specific analysis
**Implementation**: See `layered-docs-system/` for the implementation based on this research

> **Note**: This research informed the Layered Documentation System, now at v3.0 with
> appendix support, confidence-based classification, and content thresholds.

---

## Table of Contents

1. [Decision Recording Systems](#decision-recording-systems)
2. [Product/Requirements Documentation](#productrequirements-documentation)
3. [Technical Specification Systems](#technical-specification-systems)
4. [Architecture Documentation](#architecture-documentation)
5. [Behavior-Driven Documentation](#behavior-driven-documentation)
6. [Reference Documentation](#reference-documentation)
7. [Meta-Systems](#meta-systems)
8. [Documentation Toolchains](#documentation-toolchains)
9. [Integration Strategies](#integration-strategies)
10. [Preventing Drift and Contradictions](#preventing-drift-and-contradictions)
11. [Versioning and Changelog Management](#versioning-and-changelog-management)
12. [Recommendations Summary](#recommendations-summary)

---

## Decision Recording Systems

### ADR (Architectural Decision Records)

**What it is**: A document that captures an important architectural decision made along with its context and consequences. ADRs record decisions that "affect the structure, non-functional characteristics, dependencies, interfaces, or construction techniques."

**Original Format** (Michael Nygard, 2011):
- Title
- Status (Proposed → Accepted/Rejected → Superseded)
- Context (the forces at play)
- Decision (the response to these forces)
- Consequences (the resulting context after applying the decision)

**Key Benefits**:
- Early identification of design issues when changes are still cheap
- Achieving consensus around a design
- Ensuring consideration of cross-cutting concerns
- Scaling knowledge of senior engineers
- Forming organizational memory around design decisions
- Acting as a technical portfolio artifact

**Best Practices**:
- Keep ADRs focused - one decision per document
- Store in repository: `/docs/adr/` or `/adr/`
- Use Markdown format (`.md`) - lean and diff-friendly
- Numbering: `0001-title.md`, `0002-title.md` (zero-padded)
- Include timestamps for time-sensitive aspects
- Keep each ADR to one page; link to appendices for detail
- Don't delete superseded ADRs; add "Superseded by ADR-00XX" note
- Treat ADRs like code - use PR reviews, comments, approvals

**When to Create**:
- Decision impacts system qualities (availability, performance, security, cost, compliance)
- Decision is hard to reverse (data model, language/runtime, cloud selection)
- Heuristic: If the decision will be referenced six months later or could save someone a week of rediscovery

**File Organization**:
```
project/
├── docs/
│   └── adr/
│       ├── 0001-use-event-sourcing.md
│       ├── 0002-choose-postgres-over-mysql.md
│       └── 0003-adopt-microservices.md
```

### MADR (Markdown Architectural Decision Records)

**What it is**: A streamlined, structured template for ADRs that emphasizes considered options with pros/cons. Created as a lean template that "fits development style."

**Templates Available**:
- `adr-template.md` - Full version with explanations
- `adr-template-minimal.md` - Mandatory sections only with explanations
- `adr-template-bare.md` - All sections, no explanations
- `adr-template-bare-minimal.md` - Mandatory sections only, no explanations

**Structure** (beyond Nygard's format):
- Decision Drivers (what motivated this decision)
- Considered Options (alternatives evaluated)
- Pros and Cons of Options (crucial for understanding rationale)
- More Information (links, references)

**Key Advantage**: The considered options with their pros and cons are crucial to understand the reasons for choosing a particular design. This makes MADR more comprehensive than simpler ADR formats.

**Tool Support**:
- ADR Manager (VS Code) - plugin with basic and professional modes
- ADR Manager (Web) - connects to GitHub to render all ADRs
- Works with standard Markdown editors

**Latest Version**: MADR 4.0.0 (released 2024-09-17)

### Y-Statements (Lightweight ADR)

**What it is**: A lean, one-line decision capturing format from "Sustainable Architectural Decisions" by Zdun et al.

**Basic Format**:
```
In the context of <use case/user story>,
facing <concern>
we decided for <option>
to achieve <quality>,
accepting <downside>.
```

**Long Format**:
```
In the context of <use case/user story>,
facing <concern>,
we decided for <option>
and neglected <other options>,
to achieve <system qualities/desired consequences>,
accepting <downside/undesired consequences>,
because <additional rationale>.
```

**Example**:
```
In the context of the Web shop service,
facing the need to keep user session data consistent and current across shop instances,
we decided for the Database Session State pattern
and against Client Session State or Server Session State
to achieve data consistency and cloud elasticity,
accepting that a session database needs to be designed and implemented.
```

**Key Characteristics**:
- Extremely lightweight - each element is a one-line statement
- Context: functional requirement (story/use case) or architectural component
- Facing: non-functional requirement
- No sophisticated tools needed - text editor or wiki sufficient

**Tool Support**:
- Java annotation `@YStatementJustification` in e-ADR project allows capturing ADs in code
- Can be embedded directly in source files

**When to Use**: Best for teams wanting minimal overhead in decision recording, or for capturing decisions directly in code.

---

## Product/Requirements Documentation

### PRD (Product Requirements Document)

**What it is**: A document that defines the purpose, features, and behavior of a product, aligning stakeholders and guiding development. Everything needed for product completion is listed.

**Modern Approach**: Today's teams favor leaner PRDs over the lengthy documents of old. Agile PRDs focus on shared understanding, customer needs, and flexibility, avoiding overly detailed specs.

**Key Components**:
- **Overview**: Summary of product and problem it solves
- **Goals**: Objectives the product aims to achieve
- **Success Metrics**: How success will be measured
- **Target Audience**: Who the product is designed for
- **Features & Requirements**: Detailed specs and functionality
- **User Experience (UX)**: How users will interact (mockups/wireframes)
- **Assumptions & Dependencies**: Factors the product relies on
- **Out of Scope**: What won't be included
- **Key Information**: Owner, version history, project status

**Primary Benefit**: Serves as a single source of truth - clearly outlining what will be delivered in a new product or release. Ensures product, engineering, design, and go-to-market teams have shared understanding.

**Best Practices**:
- Engage stakeholders early (design, marketing, development)
- Create as collaborative effort (PMs, dev teams, UX designers, stakeholders)
- Treat as living document - update as product evolves
- Avoid overly detailed specs in Agile environments
- Ensure timeline is realistic before finalizing

**Relationship to Other Docs**:
Based on the PRD:
- Engineering creates a **Functional Specification** (how each item will be implemented)
- Engineering creates/updates an **Architectural Design Document**
- Testing creates test plans
- Marketing creates go-to-market materials

### Lean PRD Formats

**Characteristics**:
- Focus on "why" and "what" over "how"
- Minimal but sufficient detail
- Quick to write and maintain
- Emphasize user stories and outcomes

**Alternative Names**:
- One-pagers
- Product briefs
- Feature specs

### Epic-Based Documentation

**Hierarchy** (from large to small):
1. **Initiative/Theme** - Strategic objective (months/quarters)
2. **Epic** - Large body of work supporting an initiative (weeks/months)
3. **Feature** - Specific capability (weeks)
4. **User Story** - Smallest unit of work (hours/days)

**Epic Definition**:
- Any scope estimated at "weeks" or longer should be considered an epic
- More specific and measurable than initiatives
- Should show contribution to primary goal
- Consists of multiple tasks, work items, or user stories

**Epic-Based Documentation Structure**:
```markdown
# Epic: [Name]

## Goal
[Strategic objective this epic supports]

## Success Metrics
[How we'll measure success]

## User Stories
- As a [user], I want [capability] so that [benefit]
- As a [user], I want [capability] so that [benefit]

## Dependencies
[What this epic depends on]

## Out of Scope
[What this epic doesn't include]
```

**Delivery Timeframes**:
- Epics: Delivered in months
- Features: Delivered in weeks
- User Stories/PBIs: Delivered in days

### User Story Documentation

**Standard Format**:
```
As a [user type],
I want [capability],
so that [benefit].
```

**User Story Components**:
- **Front**: The story itself
- **Back**: Acceptance Criteria (AC) - agreed boundaries for successful delivery

**Acceptance Criteria with BDD**:
Stories use Behavior-Driven Development (BDD) in Given-When-Then format:
```gherkin
GIVEN [initial context]
WHEN [action occurs]
THEN [expected outcome]
```

**Best Practice**: Both the feature and individual user stories should have their own BDD scenarios testing all possible cases (positive, negative, corner cases).

**SAFe Framework Tiers** (Scaled Agile):
1. Epic
2. Capability
3. Feature
4. Story

All represent solution's intended behavior, with stories in Agile team backlogs containing detailed implementation work.

---

## Technical Specification Systems

### RFC (Request for Comments) Process

**What it is**: A process where engineering teams write design docs before starting non-trivial projects, then circulate for review and feedback. Similar to writing automated tests before code.

**Origin**: Adapted from IETF RFC process, but used internally for project design.

**Benefits**:
- Communicate intent and get feedback before development
- Give teams clear picture of what's ahead
- Find problems and improvement points before implementation
- Build solutions with high commitment from all parties
- Discussion early in design process
- More maintainable architecture through upfront review

**Companies Using RFCs**:
- **Google** - calls them "design docs"
- **Airbnb** - writes specs for both Product and Engineering
- **LinkedIn** - "strong culture of writing RFCs and doing RFC reviews"
- **Spotify** - "RFCs and ADRs are deeply embedded part of culture"
- **Uber** - scaled to 2000+ engineers using RFCs (with challenges)

**RFC Best Practices**:
- Describe scope and approach
- Focus on "what" more than "how"
- Should NOT contain specific tasks or project plan
- Pair less-experienced engineers with senior/staff engineers as 'backers'
- Use Architecture Review for big decisions (language adoption, major tech choices)
- Gather feedback early and often

**Scaling Challenges** (Uber's experience):
- **Noise**: Hundreds of RFCs weekly overwhelming experienced engineers
- **Ambiguity**: Every team had autonomy on when to write RFC, leading to inconsistency
- **Solution**: Needed clearer guidelines on which work requires RFC

**Squarespace's "Yes, if" Process**:
- Iterative RFC process
- Architecture Review used for major decisions like:
  - Adopting Go as first-class language
  - Switching to gRPC for microservices
  - Deep review of systems crossing multiple teams

### TDD (Technical Design Document)

**What it is**: A comprehensive document outlining architectural and technical aspects of a software system. Serves as roadmap for development team and communication bridge between stakeholders, designers, developers, and testers.

**Relationship to Other Docs**:
```
PRD (What system needs to do from user perspective)
  ↓
TRD (What system needs to accomplish technically)
  ↓
TDD (How system will accomplish it)
```

**Key Elements**:
- Architectural approach
- Database schemas
- Flow diagrams
- Wireframes
- Components
- Input validation
- Security considerations
- API endpoints
- Sample API requests/responses
- Pseudocode
- Alternative approaches and tradeoffs

**Best Practices**:
- Define scope clearly (what it covers AND what it doesn't)
- Identify target audience and tailor content
- If detail grows, link to appendices rather than bloating the TDD
- Create using collaborative tools (Google Docs, Confluence)
- Use diagramming tools (draw.io, Lucidchart) for visuals

**Google-Style Design Docs**:
Design docs at Google fulfill several functions:
- Early identification of design issues when changes are cheap
- Achieving consensus around design
- Ensuring consideration of cross-cutting concerns
- Scaling knowledge of senior engineers
- Forming organizational memory
- Acting as technical portfolio artifact

**Rule #1**: Write them in whatever form makes sense for the particular project (informal, not strict guidelines).

**Common Structure** (established as useful):
- Context and scope
- Goals and non-goals
- Proposed solution
- Alternative approaches considered
- System context diagram
- APIs and data models
- Considerations (security, privacy, observability, etc.)
- Milestones

---

## Architecture Documentation

### C4 Model

**What it is**: A lean graphical notation technique for modeling software architecture based on hierarchical decomposition. Created by Simon Brown (2006-2011), building on UML and 4+1 architectural view model.

**Purpose**: Help teams describe and communicate software architecture during up-front design and when retrospectively documenting existing codebases. Raises maturity of architecture diagrams.

**The Four Levels of Abstraction**:

1. **Context Diagram** (Level 1)
   - Very high-level overview
   - Shows system's relationship to external dependencies
   - Highlights external APIs, databases, third-party services
   - Shows where software sits in larger landscape
   - Audience: Everyone (technical and non-technical)

2. **Container Diagram** (Level 2)
   - Shows high-level technical building blocks
   - Containers = applications, data stores, file systems, etc.
   - For many systems, Context + Container diagrams are sufficient
   - Conveys largest part of information for most stakeholders
   - Audience: Technical stakeholders

3. **Component Diagram** (Level 3)
   - Decomposition of each container into components
   - Only needed for complex systems
   - Shows internal structure of containers
   - Audience: Developers, architects

4. **Code Diagram** (Level 4)
   - Internal structure of components (classes/objects)
   - Often auto-generated by IDEs
   - C4 doesn't define notation at this level (use UML)
   - Audience: Developers working on specific components

**Key Benefits**:
- Facilitates collaborative visual architecting
- Supports evolutionary architecture in agile teams
- Breaks down barriers between stakeholders and developers
- Provides common language for describing systems
- Interactive diagrams (zoom in/out)
- Avoids overly formal documentation

**Tools**:
- **Structurizr** - created by Simon Brown, specifically for C4 model
  - Interactive diagrams (zoom, pan)
  - Animatable, embeddable
  - Auto-generated diagram key/legend
- **C4-PlantUML** - combines PlantUML with C4 model
- **IcePanel** - collaborative diagramming with drag-and-drop UI

**When to Use Each Level**:
- Always start with Context (Level 1)
- Almost always include Container (Level 2)
- Add Component (Level 3) only for complex systems
- Use Code (Level 4) sparingly, often auto-generate

### arc42 Template

**What it is**: A comprehensive template for effective, practical software architecture documentation and communication. Created by Dr. Gernot Starke and Dr. Peter Hruschka (2005). Free and open source.

**Answers Two Questions**:
1. What should you document/communicate about your architecture?
2. How should you structure this information?

**Key Features**:
- Structured and modular approach
- Covers all relevant aspects efficiently
- Flexibly customizable for any project size/complexity
- Tool and technology agnostic
- Can be used for arbitrary systems
- Completely free to use

**Main Template Sections**:

1. **Introduction and Goals**
   - Requirements overview
   - Quality goals (top 3-5 for major stakeholders)
   - Stakeholder table with expectations

2. **Constraints**
   - Technical constraints
   - Organizational constraints
   - Conventions

3. **Context and Scope**
   - Business context
   - Technical context

4. **Solution Strategy**
   - Fundamental decisions and solution strategies
   - Technology decisions
   - Top-level decomposition
   - Approaches to achieve quality goals
   - Organizational decisions

5. **Building Block View**
   - Static decomposition
   - Hierarchy of white boxes containing black boxes
   - Source code abstractions
   - Appropriate level of detail

6. **Runtime View**
   - Behavior of building blocks
   - Runtime scenarios

7. **Deployment View**
   - Technical infrastructure
   - Environments, computers, processors, topologies
   - Mapping of building blocks to infrastructure

8. **Cross-Cutting Concepts**
   - Overall regulations and solution approaches
   - Relevant in multiple parts of system
   - Often related to multiple building blocks

9. **Architecture Decisions**
   - Important decisions and their rationale

10. **Quality Requirements**
    - Quality tree
    - Quality scenarios

11. **Risks and Technical Debt**
    - Known risks
    - Technical debt

12. **Glossary**
    - Important domain and technical terms

**Benefits**:
- Saves time with proven structure
- No need to reinvent the wheel
- Standardized documentation across teams
- Facilitates maintenance and onboarding
- Reduces misunderstandings
- Promotes efficient collaboration
- Efficient while ensuring high quality

**Tool Support**:
- Tool agnostic - use your choice:
  - Wiki
  - Office documents
  - Markdown or AsciiDoc
  - UML tools
  - LaTeX
- Extensive documentation with 140+ tips and 30+ examples

**Compatibility**:
- Can be used with C4 model (arc42 for software, C4 for higher-level IT architecture)
- Can be used with TOGAF
- Organizations often use multiple frameworks together

**When to Use**:
- Projects requiring comprehensive architecture documentation
- Teams needing standardized documentation approach
- Long-lived systems requiring maintainable docs
- Projects with multiple architects or changing team members

---

## Behavior-Driven Documentation

### Gherkin/Cucumber

**What it is**: A set of grammar rules making plain text structured enough for automated testing tools to understand. Business-readable language for describing behavior without implementation details.

**Purpose**: Gherkin feature specifications serve multiple roles:
- Living documentation that never becomes outdated
- Executable specifications of the system
- Automated tests
- High-level end-user perspective of product behavior

**Key Concept**: Specification by Example - concrete examples illustrate rules and behavior.

**Language**: Domain-specific language for defining tests in Cucumber format. Uses plain language to describe use cases, removing logic details from behavior tests.

**Keywords**:
- `Feature:` - Name and description of feature
- `Rule:` - Business rule (added 2018)
- `Scenario:` - Concrete example
- `Given` - Initial context
- `When` - Action/event
- `Then` - Expected outcome
- `And`, `But` - Additional steps

**Structure**:
```gherkin
Feature: User Authentication
  As a user
  I want to log in to the system
  So that I can access my account

  Rule: Users must provide valid credentials

    Scenario: Successful login with valid credentials
      Given I am on the login page
      When I enter valid username and password
      And I click the login button
      Then I should be redirected to my dashboard
      And I should see a welcome message

    Scenario: Failed login with invalid credentials
      Given I am on the login page
      When I enter invalid username or password
      And I click the login button
      Then I should see an error message
      And I should remain on the login page
```

**Rules and Examples** (Critical Relationship):
- Each example illustrates one (and only one) rule
- Without examples, a rule may be ambiguous
- Without a rule, an example lacks context
- Together they fully specify expected behavior

**Writing Better Gherkin**:
- Describe WHAT (intended behavior), not HOW (implementation)
- Use declarative style (behavior over implementation details)
- Declarative scenarios read better as "living documentation"
- Avoid too many steps - reduces expressive power
- Keep scenarios focused and concise

**As Documentation, Specification, and Tests**:
- **Documentation**: Living documentation always in sync with code
- **Specification**: Concrete examples of how features work
- **Tests**: Executable specifications that verify behavior

**Example Map Relationship**:
- Theme → Epic → Feature → Rule → Example
- Each level provides context for the next
- Examples validate rules
- Rules provide context for examples

**Living Documentation**:
Feature files are living documentation because:
- Updated together with code changes
- Never becomes outdated
- High-level user perspective
- Executable and verifiable
- Accessible to non-technical stakeholders

**Popular Format**: Gojko Adzic's poll showed Given/When/Then received 71% of votes as most popular format for expressing examples.

### SpecFlow+ LivingDoc

**What it is**: Tool for generating living documentation from SpecFlow/Gherkin feature files. Combines all scenarios from all feature files into consolidated HTML report.

**Problem it Solves**:
- Feature files stored in repository with automation code
- Requires Visual Studio or Git to view
- Inaccessible to non-coding stakeholders
- LivingDoc breaks down this barrier

**Key Features**:
- Always up to date with feature files and test results
- Automatically updated when feature files change or tests run
- Synchronized with latest version of test suite
- Interactive and dynamic (not static document)

**Capabilities**:
- Filter scenarios by tags, features, or results
- Search scenarios by keywords, phrases, or regex
- View summary and statistics
- Link scenarios to user stories/requirements (traceability)
- Show alignment with business goals

**Tool Options**:
- **SpecFlow+ LivingDoc CLI** - command-line tool for updates
- Manual or CI/CD pipeline integration
- Merge multiple test results from different sources
- Generate local or self-hosted HTML
- With or without test results

**Publishing Options**:
- Azure DevOps
- GitHub Pages
- Own web server

**SpecSync Integration**:
- Synchronizes BDD scenarios to Azure DevOps Test Cases
- Indicates automated tests
- Publishes execution results to Test Cases
- Complete traceability and "living documentation"

**Note**: While SpecFlow has been retired, the community continues through ShiftSync for quality engineering resources.

---

## Reference Documentation

### CLI Reference Standards

**Key Resources**:
- **clig.dev** - Modern CLI guidelines updating UNIX principles
- **GNU Coding Standards** - POSIX guidelines for CLI options
- **Google Developer Documentation Style Guide** - CLI syntax documentation

**Best Practices**:

**Flag Design**:
- Have full-length versions of all flags (both `-h` and `--help`)
- Use one-letter flags only for commonly used flags
- Use full versions in scripts for verbosity and clarity
- Don't "pollute" namespace of short flags

**Standard Options** (GNU):
- `--version` - show version information
- `--help` - show help/usage information

**Output Format**:
- Plain, line-based text - easy to pipe between commands
- JSON for structured data when needed
- Allows integration with web tools
- Supports both traditional and modern workflows

**Help System**:
- `--help` provides essential documentation
- Lets new users discover all commands and options
- Reference for experienced users
- Should provide complete list of:
  - Commands
  - Subcommands
  - Short-name equivalents
  - Simple description for each

**Help Screen Structure**:
- Getting started documentation for CLI
- List commands in logical groups
- Each group in alphabetical order
- Quick and thorough explanation for each command

**Documentation Content**:
- Break information into digestible chunks
- Support scanning by action-oriented users
- Don't overwhelm with every detail in CLI
- Users skim-read, looking for immediate needs

**Formatting**:
- Use monospace font for command-line parts
- Use normal text for descriptions and explanations
- Avoid optional argument syntax in click-to-copy commands
- Square brackets, curly braces, pipes can break commands if not removed

**Usability Features**:
- **Tab completion** - document in format for popular shells
- **Sane defaults** - choose reasonable defaults, document flags to change
- **Good error messages** - parallel best practices from other UIs
- Follow conventions users expect

**Documentation Best Practices**:
- Provide inline link to command reference
- Good place: text introducing the command
- Optional/mutually exclusive arguments can break if copied directly
- Warn users to remove syntax characters

**Content Structure**:
```
command-name - Brief description

Usage:
  command-name [OPTIONS] <REQUIRED_ARG>

Options:
  -f, --flag          Description
  -o, --option VALUE  Description

Examples:
  command-name --flag input.txt
  command-name -o value input.txt
```

### API/Configuration Documentation

**OpenAPI Specification** (for REST APIs):
- Standard, language-agnostic interface description for HTTP APIs
- Allows humans and computers to discover capabilities
- No need for source code, additional docs, or network inspection

**Format Options**:
- **YAML** (typically preferred)
  - Slightly reduced file size
  - More readable for humans
  - No commas, brackets, or quotes required
  - Relies on indentation
  - Must use YAML 1.2 for complete interchangeability
- **JSON**
  - Completely interchangeable with YAML
  - Requires commas, brackets, quotes
  - No comments support
  - Easier for machine generation

**YAML vs JSON Considerations**:
YAML Requires:
- Hyphens before array items
- Relies heavily on indentation

JSON Requires:
- Commas separating fields
- Curly brackets around objects
- Double quotes around strings
- Square brackets around arrays

**Document Structure**:
Must contain OpenAPI Object with at least:
- `openapi` - version field
- `info` - metadata
- One of: `paths`, `components`, or `webhooks`

**Recommended Naming**:
- `openapi.json`
- `openapi.yaml`

**Reusable Components**:
- Global `components` section for schemas
- Referenced in individual endpoints
- Declares schemas to be called elsewhere
- Avoids duplication across endpoints

**Versions**:
- **OpenAPI 3.1.x** - based on JSON Schema Draft 2020-12
- **OpenAPI 3.0.x** - based on JSON Schema Wright Draft 00
  - Versions 3.0.0 through 3.0.4 are functionally identical

**Beyond REST APIs**:
For general configuration documentation:
- Use JSON Schema for structure validation
- Provide examples for each configuration option
- Document default values
- Explain impact of each setting
- Group related options logically
- Provide migration guides for breaking changes

**Configuration Reference Pattern**:
```markdown
# Configuration Reference

## Option: `setting_name`

**Type**: string | number | boolean | object | array
**Default**: `default_value`
**Required**: yes | no

**Description**: What this setting does and when to use it.

**Example**:
```yaml
setting_name: example_value
```

**Related Settings**: Links to related configuration options
```

---

## Meta-Systems

### Docs-as-Code

**What it is**: A philosophy that documentation should be written with the same tools and workflows as code, integrated into the product team.

**Core Principles** (Four Key Pillars):
1. Same quality assurance, testing, and CI/CD as code
2. Version control with Git (preferably same repo as code)
3. Plain text markup formats (Markdown, reStructuredText, AsciiDoc)
4. Automated deployment, linting, and validation

**Key Benefits**:
- **Better Collaboration**: Developers, technical writers, DevOps contribute the same way
- **Improved Accuracy**: Docs alongside codebase more likely updated in sync
- **Automation**: CI/CD pipelines automate deployment, linting, validation
- **Reduced Manual Work**: Less administrative overhead

**Modern Challenges (2026)**:
- Deep system behavior poorly documented
- Process docs and style guides exist
- The underlying "why this works this way" rarely captured
- Missing engineering intent across repositories
- Loss of context required to maintain coherence at scale

**Best Practice**: Document the "why" behind decisions (the intent), not just the "what" (the implementation).

**Implementation Examples**:

**Squarespace Domains Team Approach**:
- Docs and code versioned together in Git
- Pull request review same for both
- Code and documentation changes in same PR
- Merged changes flow through CI/CD to production

**Goal**: Empower engineers to write technical documentation frequently and keep it up to date by integrating with their tools and processes.

**Mindset**: Docs-as-Code is not just tools or process, but a mindset of creating and maintaining technical documentation integrated into development workflow.

### Diátaxis Framework

**What it is**: A systematic approach to organizing technical documentation created by Daniele Procida. Name from Ancient Greek "across arrangement."

**Core Concept**: Identifies four distinct documentation needs with four corresponding forms, placed in systematic relationship.

**The Four Types**:

1. **Tutorials** (Learning-Oriented, Practical)
   - Takes student by hand through learning experience
   - Always practical - user does something under guidance
   - Instructor responsible for learner's safety and success
   - Analogy: Driving lesson
   - Reader "on rails" - specific destination with exact steps
   - Reader should follow start to finish without decision-making
   - **When**: User is studying/learning

2. **How-To Guides** (Problem-Oriented, Practical)
   - Addresses real-world goal or problem
   - Provides practical directions
   - For already-competent users
   - Concerned with work rather than study
   - Helps reader reach destination of their choosing
   - Not precise example, but general guidance
   - **When**: User is working/applying knowledge

3. **Reference** (Information-Oriented, Theoretical)
   - Technical description - facts
   - Accurate, complete, reliable information
   - Free of distraction and interpretation
   - Architecture should reflect structure of thing being described
   - Like a map or dictionary - factual, precise, structured
   - Look things up when needed, don't read cover-to-cover
   - Examples: API docs, command references
   - **When**: User needs specific information

4. **Explanation** (Understanding-Oriented, Theoretical)
   - Provides context and background
   - Helps understand and see bigger picture
   - Joins things together
   - Answers "why?" questions
   - Can contain opinions and perspectives
   - May circle around subject from different directions
   - Like reading science of cooking vs following recipe
   - **When**: User wants to deepen understanding

**The Diátaxis Grid**:

```
                    PRACTICAL
                       |
        Tutorials  |  How-To Guides
    (Learning)     |  (Goals)
                   |
STUDY  -------------------------  WORK
                   |
     Explanation   |  Reference
   (Understanding) |  (Information)
                   |
                  THEORETICAL
```

Horizontal Axis: Acquisition (study) → Application (work)
Vertical Axis: Action (practical) → Cognition (theoretical)

**Key Distinctions**:
- **Tutorials vs How-To**: Learning vs applying (both practical)
- **Explanation vs Reference**: Understanding vs information (both theoretical)
- **Tutorials vs Explanation**: Practice vs theory (both about learning)
- **How-To vs Reference**: Goals vs facts (both for work)

**Benefits**:
- Solves problems with content (what to write)
- Solves problems with style (how to write it)
- Solves problems with architecture (how to organize it)
- Serves both users and creators/maintainers
- Light-weight, easy to grasp, straightforward to apply
- No implementation constraints
- Not tied to specific tooling

**Adoption**:
- **Gatsby** - reorganized open-source documentation
- **Cloudflare** - used as "north star for information architecture"
- **Canonical/Ubuntu** - foundation for documentation

**vs DITA**: Diátaxis isn't an XML schema, doesn't require validation, not tied to specific CMS tools.

### Divio Documentation System

**What it is**: Earlier name for what became Diátaxis. Same four-part framework.

**The Four Functions**:
1. Tutorials - take user through series of steps
2. How-to guides - solve real world problems
3. Reference guides - technical explanations of how things work
4. Explanations - clarify topics and how project fits in landscape

**Key Principle**: Each requires distinct mode of writing. People need these four different kinds at different times in different circumstances.

**Structure Requirements**:
- Documentation must be explicitly structured around the four types
- All four must be kept separate and distinct
- Division makes obvious what material goes where (for authors and readers)

**Benefits**:
- Make documentation better by doing it the right way
- Right way is easier way - easier to write and maintain
- Less stressful documentation in long term

**When to Use Which**:
- Simple documents (README.md): Markdown often sufficient
- Complex requirements: Use more powerful system (AsciiDoc, full framework)
- Single page: Markdown
- Multi-page documentation with dedicated reader: Diátaxis structure

### Documentation-Driven Development (DDD)

**Philosophy**: From user perspective, if a feature is not documented, it doesn't exist. If documented incorrectly, it's broken.

**Core Approach**: Documentation written before code implementation. Serves as blueprint for development.

**Key Steps**:

1. **Document the feature first**
   - Describe to users before any code
   - If not documented, doesn't exist
   - Best way to define feature in user's eyes

2. **Review with users**
   - Documentation reviewed by users before development
   - Catches misunderstandings early

3. **Develop with TDD**
   - Test-driven development preferred
   - Unit tests test features as described by documentation
   - If functionality misaligns with docs, tests fail

**Benefits**:
- Creates clear expectations and alignment before development
- Reduces misunderstandings between stakeholders and developers
- Forces thorough thinking about design and architecture upfront
- Serves as contract between teams/services
- Facilitates better code maintenance and onboarding
- Enables parallel work through clear interfaces
- Promotes better testing through documented expectations
- Self-feedback cycle on APIs and scope of work

**Relationship to Other Methodologies**:
- **Not a replacement** for TDD, but augmentation
- Summary: document-program-test-repeat
- Complements Test-Driven Development (tests first)
- Complements Behavior-Driven Development (behavioral specs)
- Similar to API-First Development (API design before implementation)

**Application**:
- Works for greenfield projects (design-first documentation)
- Works for brownfield projects (documenting during evolution)
- Can be applied to both new features and refactoring

**Modern AI Enhancement**: AI-driven approaches now support documentation-driven workflows from documentation to deployment.

---

## Documentation Toolchains

### Markdown vs AsciiDoc

**Markdown**:

**Pros**:
- De facto standard, widely available
- Easy to learn
- Popular with technical audiences
- Used on GitHub, GitLab, etc.
- Many familiar with it
- Good for single-page documents
- Important when reading plain text matters
- Sufficient for simple documents (READMEs)

**Cons**:
- Multiple incompatible flavors (GitHub, CommonMark, MultiMarkdown)
- Subtle differences between flavors
- Inconsistent rendering across platforms
- Limited for complex needs (cross-references, footnotes, etc.)
- Requires embedded HTML for advanced features
- Tables limited (no rowspan/colspan)
- Reaches limits quickly for complex documentation

**AsciiDoc**:

**Pros**:
- Powerful markup for complex structures
- Particularly suited for technical documents and manuals
- Syntax designed for extensibility
- Maximally reusable content
- Wide range of text formatting (footnotes, cross-references, index entries)
- Tables with extensive formatting options
- Clearly defined specifications
- Consistency in document creation
- Suitable for collaborative work
- Suitable for long-term maintenance
- Multiple output formats (HTML, PDF, EPUB)
- GitHub renders AsciiDoc as well as Markdown

**Cons**:
- Less familiar to engineers and writers
- Syntactic idiosyncrasies (multiple asterisks for nested lists)
- Smaller community than Markdown
- Requires specialized toolchain

**Recommendations**:

**When to Use Markdown**:
- Single-page documents
- Simple documentation (READMEs)
- Quick notes and comments
- Reading plain text is important
- Team very familiar with Markdown

**When to Use AsciiDoc**:
- Multi-page documentation
- Complex technical documentation
- Software architecture documentation
- Dedicated reader/location acceptable
- Long-term maintainability critical
- Need for extensibility

**Best Toolchain for Engineering Docs**:
- **AsciiDoc + Asciidoctor + Antora**
- AsciiDoc: Consistent, extensible, semantically rich
- Asciidoctor: Powerful rendering, real-time previews, CI/CD integration
- Antora: Static site generator for AsciiDoc, supports modular, versioned, multi-repo documentation

**Tool Support**:
- GitHub: Renders both Markdown and AsciiDoc
- Jekyll: Seamless AsciiDoc integration
- Most modern doc platforms: Support both formats

---

## Integration Strategies

### How Documentation Systems Work Together

**Document Hierarchy and Workflow**:

```
1. PRD (Product Requirements Document)
   ↓ Defines WHAT system needs to do from user perspective

2. TRD (Technical Requirements Document)
   ↓ Translates to WHAT system needs to accomplish technically

3. TDD (Technical Design Document)
   ↓ Describes HOW system will accomplish it

4. Implementation
   ↓

5. ADR (Architectural Decision Records)
   Documents significant decisions made during implementation
```

**Three-Document Pattern**:

1. **PRD (Problem Requirements Document)** - Problem that needs solving
2. **RFC (Request for Comments)** - Proposed solution and design
3. **ADR (Architectural Decision Record)** - Decisions made

**Benefits of This Pattern**:
- PRDs and RFCs build culture of writing to make decisions
- Makes decision-making transparent and accessible
- Any decisions outside PRD/RFC captured in ADRs
- Creates iterative, cumulative history of all decisions
- If someone asks "how does your system work?", give curated list of RFCs

**Living Documentation Through Integration**:
Writing to make decisions and design generates:
- High-quality documentation
- Living documentation (stays current)
- Democratizes decision-making process
- Scales decision-making across organization

**Document Relationships**:

**From PRD, Create**:
- Functional Specification (Engineering) - how each item implemented
- Architectural Design Document (Engineering) - system architecture
- Test Plans (QA) - how to verify requirements
- Go-to-Market Materials (Marketing) - how to position product

**Common Information**:
Requirements documents contain overlapping information, making creation less time-consuming than it appears. Leverage this overlap.

### Single Source of Truth (SSOT) Strategies

**What it is**: Structuring information so every data element is mastered (edited) in only one place, then referenced elsewhere. Not a tool or system, but a state of being for data.

**Key Principles**:

1. **Linking Strategy**
   - Crucial to avoid duplicate content
   - Perfect SSOT is self-contained
   - Every data point exists once
   - Fully developed linking strategy

2. **Component-Based Approach**
   - System is component-based (not page/document/file-based)
   - Enables reuse via linking across departments/company
   - Essential for SSOT
   - Without componentization, eliminating redundancy impossible

3. **Rich Information**
   - Includes photos, support documents
   - Links everything to related files and directories
   - Nothing is missing
   - Everything documented thoroughly

4. **Accessibility**
   - Cloud-based for anywhere, anytime access
   - Maximum accessibility
   - Must be accessible to all parties depending on the truth
   - Communicates when truth changes

**Implementation Challenges**:
- Ideal SSOT rarely possible in enterprises
- Organizations have multiple information systems
- Systems need access to same data (e.g., customer)
- Often purchased COTS products that can't be easily modified
- Integration biggest hurdle - many orgs have 900+ applications

**Best Practices**:

1. **Centralization**
   - All documents in single system
   - Eliminates duplication
   - Workers know where proper version is
   - Administrators can track access/edits

2. **Version Control**
   - Track revisions over time
   - Compare different iterations
   - Revert to previous states
   - Clear history of evolution
   - Only latest approved versions accessible

3. **Standardization**
   - Unified content reuse (setup steps, definitions, integrations)
   - Predefined templates for common documents
   - Consistent structure and format
   - Reduces repetitive work

4. **Avoid Multiple Sources of Truth**
   - Determine one place to store knowledge
   - Reserve time to discuss with team
   - Block time to move information
   - Ongoing effort, not one-time project

**Visual Software Benefits**:
- Displays resources on infinite canvas
- Complete visibility into SSOT
- See critical relationships between documents

---

## Preventing Drift and Contradictions

### The Problem: Documentation Drift

**What it is**: Ongoing process of codebase becoming out of sync with documentation. Occurs when features/improvements/changes made to code without updating documentation.

**Consequences**:
- Confusion within development team
- Documentation no longer used
- Lost trust in documentation
- Wasted time searching for truth

### Challenges in Multi-Document Systems

**Multiple Products/Documents**:
- Content lives in silos
- Fragmented information architecture
- Time lost searching
- Duplicated work
- Fragmented user experiences

**Two or More Disconnected Systems**:
- Systems easily get out of sync
- Confusion on which document is correct
- Products evolve at different speeds
- Documentation becomes misaligned
- Conflicting and outdated information
- Volume makes manual consistency checks difficult and error-prone

### Prevention Strategies

**1. Integrate Documentation into Development Process**:
- Make it easy for developers to write documentation
- Make it required part of development/release process
- Add documentation directly into codebase (markdown files)
- Use tools like MkDocs to auto-generate documentation sites
- Documentation changes in same PR as code changes

**2. Version Control**:
- Robust version control for documentation
- Track revisions over time
- Compare different iterations
- Revert to previous states when necessary
- Clear history of document evolution
- Prevents using outdated information
- Only latest approved versions accessible

**3. Centralization**:
- Single system for all documents
- Eliminate numerous files across locations
- Workers no longer wonder where proper version is
- Administrators track who opened/edited documents
- Central documentation platform essential

**4. Standardization and Content Reuse**:
- Unified documentation for common content
- Reuse: setup steps, definitions, integration instructions
- Eliminates repetitive work
- Predefined templates for documents
- Consistent structure and format across all docs

**5. Automation**:
- Automate document control functions
- Auto-sync current revisions from systems (e.g., Autodesk Revit)
- Connect/automate workflows across systems
- Reduce manual administrative tasks
- Reduce delays from manual processes

**6. Audit Trails and Tracking**:
- Document control system tracks changes
- Only latest approved versions in circulation
- Detailed audit trail for all changes:
  - Who made changes
  - When changes made
  - Why changes made

**7. Clear Linking Strategy**:
- Crucial for avoiding duplicate content
- Component-based, not file-based
- Enable reuse via linking
- Make relationships between documents explicit

**8. Automated Consistency Checks**:
- Formal analysis for automatic error detection
- Type errors, nondeterminism, missing cases, circular definitions
- Consistency checking for requirements specifications
- Automated validation scripts
- Integration with validation tools and frameworks

### Resolving Existing Drift

**Two Approaches**:

1. **Eventual Consistency**:
   - When developers update feature, check that page for other issues
   - Eventually documentation becomes up to date as features worked on
   - Require developers reading docs to fix issues they find

2. **Dedicated Cleanup**:
   - Systematic review of all documentation
   - Compare with current codebase
   - Update all discrepancies
   - More time-intensive but faster results

### Automated Testing for Documentation

**Benefits**:
- Automation improves efficiency, consistency, scalability
- Automated tests execute faster than manual
- Reduces time-to-market
- Eliminates human errors common in manual testing
- Consistent execution with same parameters
- Ensures reliability and accuracy

**Types of Automated Checks**:
- Link checking (no broken links)
- Code example validation (examples actually work)
- API endpoint verification (documented endpoints exist)
- Schema validation (configs match schemas)
- Consistency checks (terms used consistently)
- Completeness checks (required sections present)

**Best Practices**:
- Deploy checks throughout pipeline
- Catch issues at each transformation step
- Clear documentation of validation rules
- Ensures rules remain consistent and understandable
- Automated code review tools for thoroughness
- Enhances efficiency, reduces human error
- Maintains high-quality standards

---

## Versioning and Changelog Management

### Semantic Versioning (SemVer)

**Format**: X.Y.Z where X, Y, Z are non-negative integers
- **X** = Major version (backward incompatible changes)
- **Y** = Minor version (backward compatible functionality)
- **Z** = Patch version (backward compatible bug fixes)

**Rules**:
- Each element must increase numerically
- Major changes break compatibility
- Minor changes and patches are backward compatible
- Minor and patches can apply to older version
- Major changes require newly released version

**Application to Documentation**:
```
1.0.0 - Initial release
1.0.1 - Fix typos, clarify examples
1.1.0 - Add new section, new examples
2.0.0 - Restructure document, breaking changes to examples
```

### Documentation Versioning Strategies

**Parallel Versioning** (for APIs):
- Maintain documentation for each version
- Software APIs change with releases
- Developers need access to both current and legacy versions
- Clear version indicators
- Easy navigation between versions

**Synchronized Versioning**:
- Keep source code and artifact versions synchronized
- Helps correlate artifact to source code state
- For multi-language docs: coordinate across all language variants
- Translation workflow integration
- Version synchronization tracking

**Independent Versioning**:
- Each document maintains its own version
- Documents evolve at different rates
- Increment versions independently
- Version mismatches are intentional and acceptable

### Changelog Best Practices

**Content**:
- List all changes in each version
- Help developers, users, stakeholders understand updates
- Include all relevant details
- Highlight breaking changes
- Provide rationale for changes
- Include impact assessment
- Include migration guidance

**Structure**:
```markdown
# Changelog

## [2.0.0] - 2026-01-15

### Breaking Changes
- Changed X to Y because [rationale]
- Migration guide: [steps]

### Added
- New feature Z
- Support for [capability]

### Changed
- Improved [feature] to [benefit]

### Fixed
- Bug in [component] causing [issue]

### Deprecated
- Feature X will be removed in 3.0.0
```

**Automated Generation**:
- Use conventional commits
- Tools: standard-version, semantic-release
- As long as git commits are conventional, get:
  - Automatic semver type
  - Free CHANGELOG generation
- Example: semantic-release/changelog plugin

**Format Standards**:
- Follow Keep a Changelog format
- Use Semantic Versioning
- Document in CHANGELOG.md
- Keep in repository root
- Update with each release

### Version Synchronization

**For Multi-Document Systems**:
- Track which document versions work together
- Compatibility matrix
- Document dependencies between docs
- Version pinning for stable combinations

**Example Compatibility Matrix**:
```markdown
| PRD | Technical Spec | API Ref | Compatible |
|-----|----------------|---------|------------|
| 2.0 | 2.1           | 1.5     | ✓          |
| 2.0 | 2.0           | 1.5     | ✓          |
| 1.5 | 2.0           | 1.4     | ✗          |
```

**CI/CD Integration**:
- Automate versioning in pipeline
- Ensure immutability of versions
- Transparency through changelogs
- Automated documentation deployment
- Automated linting and validation

**Revision History Table**:
```markdown
| Version | Date | Changes |
|---------|------|---------|
| 0.5.1   | Dec 2024 | Clarified validation timing, added error examples |
| 0.5.0   | Nov 2024 | Added new validation section |
| 0.4.0   | Oct 2024 | Initial release |
```

---

## Recommendations Summary

### Systems That Work Well Together

**For Greenfield Projects** (Design-First Documentation):

```
1. Product Layer
   ├── Lean PRD or Epic-Based Documentation
   │   Define what to build and why
   │
2. Design Layer
   ├── RFC or TDD (Technical Design Document)
   │   Design how to build it (before implementation)
   │   Use MADR for considered options and tradeoffs
   │
3. Architecture Layer
   ├── C4 Model (Context + Container diagrams minimum)
   │   Visual architecture documentation
   │   AND/OR
   ├── arc42 Template (for comprehensive architecture)
   │   Full architectural documentation
   │
4. Behavior Layer
   ├── Gherkin/Cucumber Feature Files
   │   Specification by example (before implementation)
   │   Living documentation from tests
   │
5. Decision Layer
   ├── ADR (Architectural Decision Records)
   │   Capture significant decisions as they're made
   │   Use MADR format for comprehensive options analysis
   │
6. Reference Layer (Created During/After Implementation)
   ├── CLI Reference Documentation
   ├── API Documentation (OpenAPI for REST)
   ├── Configuration Reference
   │
7. End-User Layer (Diátaxis Framework)
   ├── Tutorials (learning-oriented)
   ├── How-To Guides (task-oriented)
   ├── Reference (information-oriented)
   ├── Explanation (understanding-oriented)
```

**Integration Pattern**:
```
PRD/Epic → RFC/TDD → ADRs → Implementation → Reference Docs
    ↓           ↓         ↓          ↓
 Features → Gherkin → Tests → Living Doc
    ↓
C4/arc42 (Architecture Documentation)
```

### Best Suited for Evolving Documentation (Living Documentation)

**Highest Value**:
1. **Gherkin/Cucumber + SpecFlow LivingDoc**
   - Automatically stays in sync with code
   - Executable specifications
   - Accessible to non-technical stakeholders
   - Tests verify behavior

2. **Docs-as-Code Approach**
   - Documentation in same repository as code
   - Same PR review process
   - Automated deployment via CI/CD
   - Version controlled with Git

3. **ADRs (Architectural Decision Records)**
   - Lightweight, in repository
   - Captures "why" as decisions made
   - Part of development workflow
   - Builds organizational memory

4. **Component-Based Documentation (SSOT)**
   - Write once, reference everywhere
   - Update in one place
   - Links maintain relationships
   - Reduces duplication and drift

**Medium Value**:
5. **C4 Model**
   - Can be updated incrementally
   - Visual diagrams easier to maintain than text
   - Tools like Structurizr support code-based diagrams

6. **OpenAPI/JSON Schema**
   - Can be generated from code
   - Stays in sync through automation
   - Single source of truth for APIs/configs

**Requires Active Maintenance**:
7. **PRDs, RFCs, TDDs**
   - Snapshot in time (design phase)
   - Should be updated but often aren't
   - Mitigation: Link to ADRs for changes

8. **arc42**
   - Comprehensive but requires ongoing effort
   - Best for long-lived, stable systems
   - Update during major architecture changes

### Preventing Drift/Contradictions

**Critical Strategies**:

1. **Single Source of Truth (SSOT)**
   - Every fact mastered in exactly one place
   - Everything else links/references
   - Component-based, not file-based
   - Clear linking strategy

2. **Docs-as-Code Integration**
   - Documentation in same repo as code
   - Documentation changes in same PR
   - Reviewed together
   - Deployed together via CI/CD

3. **Automated Validation**
   - Link checking
   - Code example testing
   - Schema validation
   - Consistency checking
   - Run in CI/CD pipeline

4. **Version Control**
   - Git for all documentation
   - Track all changes
   - Audit trail of who/when/why
   - Only latest approved version accessible

5. **Clear Document Hierarchy**
   - Understand relationships between docs
   - PRD → TDD → ADRs pattern
   - Know which is authoritative for what
   - Compatibility matrix for multi-doc systems

6. **Structured Linking**
   - Explicit links between related documents
   - References instead of duplication
   - Traceability (user story → feature → ADR)
   - Visual relationship maps

7. **Regular Audits**
   - Eventual consistency approach (fix when touched)
   - Or dedicated review cycles
   - Automated tooling for drift detection
   - Require readers to fix issues they find

**Technical Measures**:

8. **Formal Consistency Checking**
   - Automatic detection of contradictions
   - Type errors, missing cases
   - Requires structured formats
   - Tools for requirements validation

9. **Synchronized Versioning**
   - Keep related docs versioned together
   - Compatibility matrices
   - Clear version indicators
   - Changelog for each version

10. **Living Documentation Priority**
    - Gherkin specs as executable tests
    - Auto-generated API docs
    - Code-based architecture diagrams
    - Minimize manually-maintained docs

### Common Tools and Formats

**Markup Languages**:
- **Markdown** - Simple docs, READMEs, quick notes, single pages
- **AsciiDoc** - Complex docs, technical manuals, multi-page, architecture

**Version Control**:
- **Git** - Universal for docs-as-code
- Store in repository with code
- Branch, merge, review like code

**Static Site Generators**:
- **MkDocs** - Markdown to documentation sites
- **Asciidoctor + Antora** - AsciiDoc to documentation sites (best for complex docs)
- **Docusaurus** - React-based, good for API docs
- **Jekyll** - GitHub Pages integration

**Diagramming**:
- **Structurizr** - C4 model diagrams as code
- **PlantUML** - Text-based UML and architecture diagrams
- **Mermaid** - Diagrams in Markdown
- **draw.io / Lucidchart** - Visual diagramming

**API Documentation**:
- **OpenAPI/Swagger** - REST API specification and docs
- **JSON Schema** - Configuration and data validation

**BDD/Living Documentation**:
- **Cucumber/Gherkin** - Behavior specs
- **SpecFlow + LivingDoc** - .NET BDD with HTML reports

**Documentation Platforms**:
- **Confluence** - Wiki-style, enterprise collaboration
- **Notion** - Flexible, component-based
- **Document360** - Technical documentation platform
- **Read the Docs** - Open source documentation hosting

**ADR Tools**:
- **adr-tools** - Command-line tools for ADR management
- **ADR Manager (VS Code)** - Plugin for managing ADRs
- **Log4brains** - ADR management with web UI

**Validation/Testing**:
- **markdownlint** - Markdown style checking
- **vale** - Prose linting
- **linkchecker** - Broken link detection
- **newman** - API endpoint validation (Postman collections)

**CI/CD Integration**:
- **GitHub Actions** - Automation for docs deployment, validation
- **GitLab CI/CD** - Similar automation
- **Netlify** - Automated deployment for static sites

### Evolution Path for a Greenfield Project

**Phase 1: Design (Before Code)**
1. Write Lean PRD or Epic breakdown
2. Create C4 Context + Container diagrams
3. Write RFC/TDD for technical approach
4. Write Gherkin feature files (specification by example)
5. Begin ADRs for major decisions

**Phase 2: Implementation**
6. Implement features with Gherkin specs as tests
7. Continue ADRs for decisions during implementation
8. Generate API docs from code (OpenAPI)
9. Create CLI reference documentation

**Phase 3: Launch**
10. Organize end-user docs with Diátaxis framework:
    - Tutorials for onboarding
    - How-Tos for common tasks
    - Reference for lookup
    - Explanation for deep understanding
11. Set up SpecFlow LivingDoc for stakeholder reports

**Phase 4: Maintenance (Living Documentation)**
12. Gherkin specs stay in sync as tests
13. ADRs continue for all significant decisions
14. API docs regenerated from code
15. C4 diagrams updated as architecture evolves
16. PRD/RFC become historical reference (link to ADRs for changes)

**Continuous Throughout**:
- All docs in Git with code
- Docs reviewed in PRs
- Automated deployment
- Automated validation
- Clear linking between documents

---

## Key Insights

### What Makes Documentation "Living"

1. **Generated from Code**:
   - API docs from OpenAPI specs
   - Architecture diagrams from DSL (Structurizr)
   - Test reports from Gherkin execution

2. **In Same Workflow as Code**:
   - Same repository
   - Same PR process
   - Same CI/CD pipeline
   - Same review standards

3. **Executable/Verifiable**:
   - Gherkin specs are tests
   - Code examples can be tested
   - API endpoints can be verified
   - Links can be checked

4. **Automated Updates**:
   - SpecFlow LivingDoc regenerated
   - CI/CD deploys changes
   - Changelogs auto-generated
   - Version numbers auto-incremented

### What Makes Documentation "Design-First"

1. **Written Before Implementation**:
   - PRDs before features
   - RFCs before code
   - Gherkin before implementation
   - TDDs before building

2. **Drives Implementation**:
   - Gherkin specs become acceptance tests
   - API specs become contract tests
   - Design docs guide architecture
   - PRDs guide features

3. **Creates Shared Understanding**:
   - Reviewed before coding
   - Stakeholder input early
   - Consensus on approach
   - Clear success criteria

### What Makes Documentation "Reference"

1. **Comprehensive Coverage**:
   - Every CLI command
   - Every API endpoint
   - Every configuration option
   - Every architectural component

2. **Structured for Lookup**:
   - Alphabetical or logical grouping
   - Clear navigation
   - Search friendly
   - Consistent format

3. **Factual and Precise**:
   - No opinions
   - No tutorials
   - Just the facts
   - Like a dictionary/map

---

## Sources

### Decision Recording Systems
- [GitHub - joelparkerhenderson/architecture-decision-record](https://github.com/joelparkerhenderson/architecture-decision-record)
- [ADR Templates | Architectural Decision Records](https://adr.github.io/adr-templates/)
- [Architectural Decision Records (ADRs) | Architectural Decision Records](https://adr.github.io/)
- [ADR process - AWS Prescriptive Guidance](https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html)
- [8 best practices for creating architecture decision records | TechTarget](https://www.techtarget.com/searchapparchitecture/tip/4-best-practices-for-creating-architecture-decision-records)
- [GitHub - adr/madr: Markdown Architectural Decision Records](https://github.com/adr/madr)
- [About MADR | MADR](https://adr.github.io/madr/)
- [The Markdown ADR (MADR) Template Explained and Distilled](https://ozimmer.ch/practices/2022/11/22/MADRTemplatePrimer.html)
- [Architecture Decision Record Template: Y-Statements | ZIO's Blog](https://medium.com/olzzio/y-statements-10eb07b5a177)
- [ADR = Any Decision Record? Architecture, Design and Beyond](https://ozimmer.ch/practices/2021/04/23/AnyDecisionRecords.html)

### Product/Requirements Documentation
- [The Only PRD Template You Need (with Example)](https://productschool.com/blog/product-strategy/product-template-requirements-document-prd)
- [What is a Product Requirements Document (PRD)? | Atlassian](https://www.atlassian.com/agile/product-management/requirements)
- [PRD Templates: What To Include for Success](https://www.aha.io/roadmapping/guide/requirements-management/what-is-a-good-product-requirements-document-template)
- [Product Requirements Document: PRD Templates and Examples](https://www.altexsoft.com/blog/product-requirements-document/)
- [Epics, Stories, and Initiatives | Atlassian](https://www.atlassian.com/agile/project-management/epics-stories-themes)
- [Story - Scaled Agile Framework](https://framework.scaledagile.com/story)
- [22/36 : BDD, Theme, Epic, Feature & Story](https://productmindset.substack.com/p/2236-bdd-theme-epic-feature-and-story)

### Technical Specification Systems
- [Software Engineering RFC and Design Doc Examples and Templates](https://newsletter.pragmaticengineer.com/p/software-engineering-rfc-and-design)
- [Companies Using RFCs or Design Docs and Examples of These - The Pragmatic Engineer](https://blog.pragmaticengineer.com/rfcs-and-design-docs/)
- [Scaling Engineering Teams via RFCs: Writing Things Down - The Pragmatic Engineer](https://blog.pragmaticengineer.com/scaling-engineering-teams-via-writing-things-down-rfcs/)
- [A Structured RFC Process](https://philcalcado.com/2018/11/19/a_structured_rfc_process.html)
- [The Power of "Yes, if": Iterating on our RFC Process](https://engineering.squarespace.com/blog/2019/the-power-of-yes-if)
- [Writing Technical Design Docs. Engineering Insights | Medium](https://medium.com/machine-words/writing-technical-design-docs-71f446e42f2e)
- [Design Docs at Google](https://www.industrialempathy.com/posts/design-docs-at-google/)

### Architecture Documentation
- [Home | C4 model](https://c4model.com/)
- [What is C4 Model? Complete Guide for Software Architecture](https://miro.com/diagramming/c4-model-for-software-architecture/)
- [The C4 Model for Software Architecture - InfoQ](https://www.infoq.com/articles/C4-architecture-model/)
- [Structurizr](https://structurizr.com/)
- [arc42 Template Overview - arc42](https://arc42.org/overview)
- [GitHub - arc42/arc42-template](https://github.com/arc42/arc42-template)
- [arc42 for your software architecture: The best choice for sustainable documentation - DEV Community](https://dev.to/florianlenz/arc42-for-your-software-architecture-the-best-choice-for-sustainable-documentation-383p)
- [Documenting software architecture with arc42 – INNOQ](https://www.innoq.com/en/blog/2022/08/brief-introduction-to-arc42/)

### Behavior-Driven Documentation
- [Reference | Cucumber](https://cucumber.io/docs/gherkin/reference/)
- [Introduction - Cucumber Documentation](https://cucumber.io/docs/guides/overview/)
- [Writing better Gherkin | Cucumber](https://cucumber.io/docs/bdd/better-gherkin/)
- [Gherkin | Cucumber](https://cucumber.io/docs/gherkin/)
- [SpecFlow+ LivingDoc - View automated test results](https://specflow.org/tools/living-doc/)
- [How do you use SpecFlow living documentation](https://www.linkedin.com/advice/3/how-do-you-use-specflow-living-documentation)
- [Improving Teamwork with SpecFlow+ LivingDoc | Automation Panda](https://automationpanda.com/2021/02/09/improving-teamwork-with-specflow-livingdoc/)

### Reference Documentation
- [Command Line Interface Guidelines](https://clig.dev/)
- [Document command-line syntax | Google developer documentation style guide](https://developers.google.com/style/code-syntax)
- [GitHub - cli-guidelines/cli-guidelines](https://github.com/cli-guidelines/cli-guidelines)
- [Command-Line Interfaces (GNU Coding Standards)](https://www.gnu.org/prep/standards/html_node/Command_002dLine-Interfaces.html)
- [10 design principles for delightful CLIs - Work Life by Atlassian](https://www.atlassian.com/blog/it-teams/10-design-principles-for-delightful-clis)
- [OpenAPI Specification v3.1.0](https://spec.openapis.org/oas/v3.1.0.html)
- [Basic Structure | Swagger Docs](https://swagger.io/docs/specification/v3_0/basic-structure/)

### Meta-Systems
- [Diátaxis framework](https://diataxis.fr/)
- [Start here - Diátaxis in five minutes](https://diataxis.fr/start-here/)
- [What is Diátaxis and should you be using it](https://idratherbewriting.com/blog/what-is-diataxis-documentation-framework)
- [How to Structure Documentation using the Diataxis Framework | Medium](https://medium.com/@techwritershub/how-to-structure-documentation-using-the-diataxis-framework-70d4a5a61db7)
- [About | Divio Documentation](https://docs.divio.com/documentation-system/)
- [Introduction | Divio Documentation](https://docs.divio.com/documentation-system/introduction/)
- [Docs as Code — Write the Docs](https://www.writethedocs.org/guide/docs-as-code/)
- [What is Docs as Code? Guide to Modern Technical Documentation | Kong Inc.](https://konghq.com/blog/learning-center/what-is-docs-as-code)
- [Making Documentation Simpler and Practical: Our Docs-as-Code Journey](https://engineering.squarespace.com/blog/2025/making-documentation-simpler-and-practical-our-docs-as-code-journey)
- [Documentation-Driven Development (DDD) · GitHub](https://gist.github.com/zsup/9434452)
- [What's Documentation-Driven Development? | Medium](https://buildwithandrew.medium.com/whats-documentation-driven-development-4b007f4de6a1)
- [A Better Way To Code: Documentation Driven Development | Playful Programming](https://playfulprogramming.com/posts/documentation-driven-development)

### Documentation Toolchains
- [Compare AsciiDoc to Markdown | Asciidoctor Docs](https://docs.asciidoctor.org/asciidoc/latest/asciidoc-vs-markdown/)
- [AsciiDoc over Markdown](https://blog.frankel.ch/asciidoc-over-markdown/)
- [Markdown or AsciiDoc: The crucial question for Docs-as-Code - embarc](https://www.embarc.de/en/markdown-vs-asciidoc/)
- [AsciiDoc: The Superior Choice for Software Engineering Documentation | Medium](https://medium.com/@nagendra.raja/%EF%B8%8F-asciidoc-the-superior-choice-for-software-engineering-documentation-ccc6f5b554db)
- [Try AsciiDoc instead of Markdown | Opensource.com](https://opensource.com/article/22/8/drop-markdown-asciidoc)

### Integration Strategies
- [Documenting Design Decisions using RFCs and ADRs - Bruno Scheufler](https://brunoscheufler.com/blog/2020-07-04-documenting-design-decisions-using-rfcs-and-adrs)
- [Decision-making and design in growing engineering organisations](https://www.form3.tech/blog/engineering/blog-decision-making)
- [Building a true Single Source of Truth (SSoT) for your team](https://www.atlassian.com/work-management/knowledge-sharing/documentation/building-a-single-source-of-truth-ssot-for-your-team)
- [Single Source of Truth Guide | Lucid](https://lucid.co/blog/single-source-of-truth-guide)
- [Single source of truth - Wikipedia](https://en.wikipedia.org/wiki/Single_source_of_truth)
- [What is a Single Source of Truth (SSOT)?](https://www.heretto.com/blog/single-source-of-truth)

### Preventing Drift and Contradictions
- [What is Documentation Drift and How to Avoid It?](https://gaudion.dev/blog/documentation-drift)
- [Multi-Product Documentation Strategy: Best Practices & Structure](https://document360.com/blog/multi-product-documentation-strategy/)
- [How to Solve 5 Common Document Control Procedure Problems](https://www.qualio.com/blog/document-control-procedure-problems)
- [Document Control: Best Practices, Compliance & Systems Guide (2025)](https://docsvault.com/blog/document-control-guide-best-practices-compliance/)
- [Automated consistency checking of requirements specifications](https://dl.acm.org/doi/10.1145/234426.234431)
- [Data Quality Testing: Key Techniques & Best Practices [2025]](https://atlan.com/data-quality-testing/)

### Versioning and Changelog Management
- [Documentation Versioning: Definition, Examples & Best Practices (2025)](https://www.docsie.io/blog/glossary/documentation-versioning/)
- [How Changelog Versioning Works (and Why It Matters)](https://announcekit.app/blog/changelog-versioning/)
- [Software versioning - Wikipedia](https://en.wikipedia.org/wiki/Software_versioning)
- [Semantic Versioning 2.0.0](https://semver.org/)
- [Using Semantic Versioning to Simplify Release Management](https://aws.amazon.com/blogs/devops/using-semantic-versioning-to-simplify-release-management/)
- [GitHub - conventional-changelog/standard-version](https://github.com/conventional-changelog/standard-version)
- [Mastering Versioning: Essential Strategies for Software Evolution](https://ones.com/blog/mastering-versioning-strategies-software-evolution/)
