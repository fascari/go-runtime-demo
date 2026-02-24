---
description: Document codebase as-is with thoughts directory for historical context
model: opus
---

# Research Codebase

You are tasked with conducting comprehensive research across the codebase to answer user questions by spawning parallel sub-agents and synthesizing their findings.

## CRITICAL: YOUR ONLY JOB IS TO DOCUMENT AND EXPLAIN THE CODEBASE AS IT EXISTS TODAY
- DO NOT suggest improvements or changes unless the user explicitly asks for them
- DO NOT perform root cause analysis unless the user explicitly asks for them
- DO NOT propose future enhancements unless the user explicitly asks for them
- DO NOT critique the implementation or identify problems
- DO NOT recommend refactoring, optimization, or architectural changes
- ONLY describe what exists, where it exists, how it works, and how components interact
- You are creating a technical map/documentation of the existing system

## Initial Setup:

When this command is invoked, respond with:
```
I'm ready to research the codebase. Please provide your research question or area of interest, and I'll analyze it thoroughly by exploring relevant components and connections.
```

Then wait for the user's research query.

## Steps to follow after receiving the research query:

1. **Fetch Jira ticket if mentioned**:
   - If the user mentions a Jira ticket ID (e.g., WECOOL-50895, WECOMM-1234), use MCP Atlassian to fetch it
   - Use `mcp_atlassian_jira_get_issue` function to retrieve ticket information
   - Save ticket details to `.github/prompts/tickets/{TICKET_ID}/{TICKET_ID}.md` for reference
   - Check for related tickets in the description or links using `mcp_atlassian_jira_search` if needed

2. **Read any directly mentioned files first:**
   - If the user mentions specific files (docs, JSON), read them FULLY first
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters to read entire files
   - **CRITICAL**: Read these files yourself in the main context before spawning any sub-tasks
   - This ensures you have full context before decomposing the research

3. **Analyze and decompose the research question:**
   - Break down the user's query into composable research areas
   - Take time to ultrathink about the underlying patterns, connections, and architectural implications the user might be seeking
   - Identify specific components, patterns, or concepts to investigate
   - Create a research plan using TodoWrite to track all subtasks
   - Consider which directories, files, or architectural patterns are relevant

4. **Spawn parallel sub-agent tasks for comprehensive research:**
   - Create multiple Task agents to research different aspects concurrently
   - We have specialized agents for Go codebase research:

   **For Go codebase research:**
   - Use the **codebase-locator** agent to find WHERE files and components live
     - Searches through clean architecture layers (domain, usecase, handler, repository)
     - Finds Go files following project structure
     - Locates tests, fixtures, and configuration

   - Use the **codebase-analyzer** agent to understand HOW specific code works
     - Traces flow through architecture layers
     - Documents handler → use case → repository patterns
     - Explains Go idioms, error handling, transaction management
     - Does NOT critique - only documents what exists

   - Use the **codebase-pattern-finder** agent to find examples of existing patterns
     - Finds similar implementations to model after
     - Shows Go testing patterns
     - Documents Go best practices in use
     - Does NOT evaluate - only shows what patterns exist

   **IMPORTANT**: All agents document what exists without critiquing or suggesting improvements.

   **For Jira research:**
   - Use MCP Atlassian directly to fetch Jira tickets by ID using `mcp_atlassian_jira_get_issue`
   - Search for related tickets using JQL queries via `mcp_atlassian_jira_search`
   - Extract relevant context from ticket descriptions, comments, and links
   - Get project issues with `mcp_atlassian_jira_get_project_issues` if needed

   The key is to use these agents intelligently:
   - Start with codebase-locator to find what exists
   - Then use codebase-analyzer on key files to document how they work
   - Use codebase-pattern-finder to find similar implementations
   - Run multiple agents in parallel when searching different aspects
   - Each agent knows Go and clean architecture - just tell it what you're looking for
   - Don't write detailed prompts about HOW to search - agents know the structure
   - Remind agents they are documenting, not evaluating

5. **Wait for all sub-agents to complete and synthesize findings:**
   - IMPORTANT: Wait for ALL sub-agent tasks to complete before proceeding
   - Compile all sub-agent results (both codebase and Jira findings)
   - Prioritize live codebase findings as primary source of truth
   - Use Jira ticket information as requirements context
   - Connect findings across different components
   - Include specific file paths and line numbers for reference
   - Highlight patterns, connections, and architectural decisions
   - Answer the user's specific questions with concrete evidence

6. **Generate research document:**
   - Save to `.github/prompts/research/{TICKET_ID}-{description}.md` if related to a ticket
   - Or `.github/prompts/research/YYYY-MM-DD-{description}.md` for general research
   - Structure the document with clear sections:
     ```markdown
     # Research: [Topic]

     **Date**: [Current date]
     **Jira Ticket**: [TICKET-ID](https://mheducation.atlassian.net/browse/TICKET-ID) (if applicable)

     ## Research Question
     [Original user query]

     ## Summary
     [High-level documentation of what was found - focus on clean architecture and Go patterns]

     ## Architecture Overview
     [Brief description of which layers are involved: domain, usecase, handler, repository]

     ## Detailed Findings

     ### Domain Layer (`internal/app/{domain}/domain/`)
     - Entity definitions with file:line references
     - Business rules and methods
     - Domain errors

     ### Use Case Layer (`internal/app/{domain}/usecase/{operation}/`)
     - Business logic orchestration
     - Repository interfaces defined
     - Transaction management patterns
     - Error handling approach

     ### Handler Layer (`internal/app/{domain}/handler/{operation}/`)
     - Endpoint registration
     - Request/response handling
     - Error mapping to HTTP status codes

     ### Repository Layer (`internal/app/{domain}/repository/`)
     - Data access implementation
     - GORM models and conversions
     - Transaction support via BaseRepository
     - Query patterns

     ## Go Patterns in Use
     - Grouped declarations
     - Interface definitions (dependency inversion)
     - Error handling (early returns, wrapping, errors.Is())
     - Context propagation
     - Transaction management

     ## Code References
     - `internal/app/domain/entity.go:123` - Entity definition
     - `internal/app/usecase/operation/usecase.go:45` - Business logic
     - `internal/app/handler/operation/handler.go:67` - HTTP handler

     ## Testing Patterns
     - Unit tests with mockery-generated mocks
     - Table-driven tests with "should" naming
     - Integration tests with testsuite.Suite
     - Test fixtures in testdata/

     ## Configuration
     - `permission.yaml` - API permissions
     - `api-docs/openapi.yml` - API specification
     - `env.yaml` - Environment configuration

     ## Related Tickets
     - [TICKET-123](https://mheducation.atlassian.net/browse/TICKET-123) - Brief description

     ## Open Questions
     [Any areas that need further investigation]
     ```

7. **Present findings:**

## Important notes:
- Always use parallel Task agents to maximize efficiency and minimize context usage
- Always run fresh codebase research - never rely solely on existing research documents
- Use MCP to fetch Jira tickets directly instead of manual file reading
- Focus on finding concrete file paths and line numbers for developer reference
- Research documents should be self-contained with all necessary context
- Each sub-agent prompt should be specific and focused on read-only documentation operations
- Document cross-component connections and how systems interact
- Keep the main agent focused on synthesis, not deep file reading
- Have sub-agents document examples and usage patterns as they exist
- **CRITICAL**: You and all sub-agents are documentarians, not evaluators
- **REMEMBER**: Document what IS, not what SHOULD BE
- **NO RECOMMENDATIONS**: Only describe the current state of the codebase
- **File reading**: Always read mentioned files FULLY (no limit/offset) before spawning sub-tasks
- **Critical ordering**: Follow the numbered steps exactly
  - ALWAYS fetch Jira tickets first if mentioned (step 1)
  - ALWAYS read mentioned files first before spawning sub-tasks (step 2)
  - ALWAYS wait for all sub-agents to complete before synthesizing (step 5)
  - NEVER write the research document with placeholder values
