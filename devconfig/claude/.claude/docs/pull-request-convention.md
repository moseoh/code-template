# Pull Request Convention

Follow these guidelines when creating pull requests. Write all content in Korean by default unless otherwise specified.

## Title Format

```
<type>[optional scope]: <description>
```

Use the same types as commit convention:

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes
- **refactor**: Code refactoring
- **test**: Test changes
- **chore**: Build/tool changes
- **perf**: Performance improvements

## Body Template

```markdown
**Summary**

- Brief description of what this PR does
- Why this change is needed
- Key benefits or impact

**Changes**

- List of specific changes made
- Technical details of implementation
- Modified files or components

**Test Plan** (optional)

- How to test the changes
- Test cases covered
- Manual testing steps

**Breaking Changes** (if applicable)

- List any breaking changes
- Migration guide if needed

**Related Issues** (if applicable)

- Closes #123
- Fixes #456
- Related to #789
```

## Examples

### Feature PR

```
feat: 사용자 프로필 이미지 업로드 기능

**Summary**
* 사용자가 프로필 이미지를 업로드하고 수정할 수 있는 기능 추가
* S3 스토리지를 통한 안전한 파일 업로드
* 이미지 크기 자동 최적화 및 썸네일 생성

**Changes**
* `/api/upload` 엔드포인트 추가
* ProfileImage 컴포넌트 구현
* 이미지 업로드 UI 및 미리보기 기능
* S3 업로드 유틸리티 함수 작성

**Test Plan**
* 다양한 이미지 포맷 업로드 테스트
* 파일 크기 제한 검증
* 프로필 이미지 변경 및 삭제 기능 확인
```

### Bug Fix PR

```
fix(auth): 로그인 세션 만료 처리 오류

**Summary**
* 세션 만료 시 무한 리다이렉트 발생하는 버그 수정
* 사용자 경험 개선을 위한 적절한 에러 메시지 표시

**Changes**
* AuthMiddleware에서 세션 만료 감지 로직 개선
* 토큰 갱신 실패 시 로그아웃 처리
* 에러 토스트 메시지 추가

**Test Plan**
* 세션 만료 상황 시뮬레이션
* 자동 로그아웃 및 리다이렉트 확인
* 에러 메시지 표시 검증

**Related Issues**
* Fixes #234
```

## Rules

- Keep title under 72 characters
- Use imperative mood in title
- Start title with lowercase letter (after type)
- Use bullet points with `*` in body sections
- Include relevant sections only (skip empty sections)
- Write in Korean unless specified otherwise
- Reference related issues when applicable
- Add reviewers and appropriate labels
- Never add "Generated with Claude Code" or similar AI attribution in commit messages or PR descriptions
- Do not include "Co-Authored-By: Claude" in commit messages
