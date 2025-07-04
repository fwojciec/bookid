# Issue Management for Solo Dev + LLM Workflow

## Overview

This document defines a minimal, efficient issue management system optimized for a solo developer using LLM agents (like Claude Code) to write all code. The system prioritizes clarity, context, and automation over complex workflows.

## Label System

### Stage Labels (Required)
Track where each issue is in your workflow:

- `stage:idea` - Rough thoughts, needs refinement before LLM can work on it
- `stage:ready` - Clear requirements, LLM can start implementation

## Issue Templates

### Standard Issue Format
```markdown
## What needs to be done
[Clear, specific description of the desired outcome]

## Context
- Related files: [List specific files if relevant]
- Dependencies: [Other issues or external factors]
- Constraints: [Technical or business constraints]

## Success criteria
- [ ] Specific, testable outcome 1
- [ ] Specific, testable outcome 2
- [ ] Tests pass
- [ ] No linting errors
```

### Quick Idea Capture
For `stage:idea` issues, just write a title. You'll add details later when moving to `stage:ready`.

## Workflow

### 1. Capture Ideas
Create issues with minimal info and `stage:idea` label. Don't overthink it.

### 2. Refine for LLM
When ready to implement:
- Add context and success criteria
- Move from `stage:idea` to `stage:ready`

### 3. LLM Implementation
When Claude Code starts work:
- Create feature branch from issue
- Implement solution
- Run tests
- Create PR (links back to issue)

### 4. Completion
When PR is created:
- Human reviews and merges
- Close issue after merge

## Best Practices

### For the Human (You)

1. **Brain dump freely** - Create `stage:idea` issues whenever you think of something
2. **Batch refinement** - Set aside time to move ideas to ready state
3. **Clear success criteria** - LLMs work best with specific, testable goals
4. **Reference files** - Always mention specific files when relevant

### For the LLM

1. **Check labels first** - Ensure issue is `stage:ready` before starting
2. **Create atomic PRs** - One issue = one PR
3. **Run all checks** - Always run tests and linting before creating PR
4. **Reference issue** - Always link PR to issue with "Fixes #X"



## Multi-Issue Projects

For projects spanning multiple issues, use a parent tracking issue:

### Parent Issue Format
```markdown
# [Project Name]

## Overview
Brief description of what we're building

## Tasks
- [ ] #45 - Task 1 description
- [ ] #46 - Task 2 description
- [ ] #47 - Task 3 description

## Success Criteria
- [ ] All sub-tasks complete
- [ ] Integration tests pass
```

### Workflow
1. Create parent issue with `stage:idea` (stays at idea - it's just tracking)
2. Create individual task issues as `stage:ready` when defined
3. Update parent issue checklist with issue numbers
4. Pin parent issue while working on the project
5. Check off tasks as their PRs are merged

The parent issue serves as a lightweight project dashboard without complex tooling.

## Example Issue Lifecycle

```
1. Human: "Need to add pagination to book search"
   → Creates issue #45 with `stage:idea`

2. Human (later): "Let me add details to #45"
   → Adds context, success criteria
   → Changes label to `stage:ready`

3. LLM: "Work on issue #45"
   → Implements pagination
   → Creates PR #46 "Fixes #45"

4. Human: Reviews PR, merges
   → GitHub automatically closes issue #45
```
