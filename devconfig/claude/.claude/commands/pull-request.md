Pull Request Convention Summary
Title Format: <type>[scope]: <description>

Same types as commit convention: feat, fix, docs, style, refactor, test, chore, perf
Keep under 72 characters, imperative mood, lowercase after type

Body Structure (include relevant sections only):

```markdown
**Summary**

- Brief description and purpose
- Why needed and key benefits

**Changes**

- Specific implementation details
- Modified components/files

**Test Plan** (optional)

- Testing approach and verification steps

**Breaking Changes** (if applicable)

- List breaking changes and migration notes

**Related Issues** (if applicable)

- Closes #123, Fixes #456, Related to #789
```

Rules:

- Write in Korean by default
- Use `*` for bullet points
- Skip empty sections
- Reference related issues
- Never include AI attribution
- Use GitHub CLI (gh) for creating PRs

Example:

```markdown
feat: 사용자 프로필 이미지 업로드 기능

**Summary**

- 사용자가 프로필 이미지를 업로드하고 수정할 수 있는 기능 추가

**Changes**

- ProfileImage 컴포넌트 구현
- `/api/upload` 엔드포인트 추가
```
