# Commit Convention

Follow Conventional Commits 1.0.0 specification for all commit messages. Write commit descriptions in Korean by default unless otherwise specified.

## Basic Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Required Types

- **feat**: Introduces a new feature (MINOR version bump)
- **fix**: Patches a bug (PATCH version bump)
- **BREAKING CHANGE**: Introduces breaking API changes (MAJOR version bump)
  - Add `!` after type (`feat!:`, `fix!:`) or
  - Use `BREAKING CHANGE:` in footer

## Optional Types

- **docs**: Documentation changes
- **style**: Code style changes (formatting, semicolons, etc.)
- **refactor**: Code refactoring without functionality changes
- **test**: Adding or modifying tests
- **chore**: Build process or tool changes
- **perf**: Performance improvements
- **ci**: CI configuration changes
- **build**: Build system changes

## Scope (Optional)

Specify module/component: `feat(parser):`, `fix(api):`, `docs(readme):`

## Examples

```
feat: 사용자 인증 기능 추가
fix(api): 로그인 타임아웃 문제 해결
docs: 설치 가이드 업데이트
feat!: API 응답 형식 변경
chore: 의존성 업데이트
refactor(utils): 날짜 포맷팅 로직 단순화
```

## Rules

- Type must be followed by colon and space (`: `)
- Use imperative mood for descriptions
- Start with lowercase letter
- No period at the end
- Keep under 50 characters when possible
- Write descriptions in Korean unless specified otherwise
- Never add "Generated with Claude Code" or similar AI attribution in commit messages or PR descriptions
- Do not include "Co-Authored-By: Claude" in commit messages
