# 자동 썸네일 생성기

assets 폴더의 로고 파일들을 자동으로 스캔하여 썸네일을 일괄 생성하는 도구입니다.

## 사용법

### 1. 로고 파일 준비

`assets/` 폴더에 다음 형식으로 로고 파일을 배치:

```
{로고이름}-{로고색상}-{배경색상}.{확장자}
```

**예시:**

- `nextjs-ffffff-000000.svg` → 흰색 로고, 검은색 배경
- `react-61dafb-282c34.svg` → 파란색 로고, 어두운 배경
- `logo-red-white.png` → 빨간색 로고, 흰색 배경

**SVG 색상 자동 변경**: SVG 파일의 경우 로고 색상이 자동으로 변경됩니다!

### 2. 실행

#### 기본 실행 (config.json 사용)

```bash
npm start
```

#### 커스텀 config 파일 사용

```bash
# 전용 스크립트 사용 (추천)
npm run bookreport

# 직접 node 실행
node generate.js -c ./config-bookreport.json
node generate.js --config ./config-custom.json

# 도움말 확인
node generate.js --help
```

#### CLI 옵션

- `-c, --config <path>`: 사용할 config 파일 경로 지정 (기본값: `./config.json`)
- `-h, --help`: 도움말 표시

### 3. 결과

- `output/` 폴더에 PNG 썸네일 생성
- 타겟 경로가 설정되어 있으면 자동 복사

## 설정 파일

### 기본 설정 (config.json)

프로젝트의 기본 설정 파일입니다.

```json
{
  "thumbnail": {
    "width": 1600,
    "ratio": "16:10",
    "defaultBackground": "#ffffff"
  },
  "logo": {
    "heightRatio": 0.4,
    "widthRatio": 0.7,
    "defaultColor": "#000000"
  },
  "paths": {
    "assets": "./assets",
    "output": "./output",
    "target": "/path/to/target/directory"
  },
  "fileConvention": {
    "separator": "-",
    "description": "파일명 형식: {로고이름}-{로고색상}-{배경색상}.{확장자}"
  }
}
```

### 커스텀 설정 파일

용도별로 다른 설정이 필요한 경우 커스텀 config 파일을 생성할 수 있습니다.

**예시: config-bookreport.json (책 표지용)**

```json
{
  "logo": {
    "heightRatio": 0.7,  // 책표지는 로고를 더 크게
    "widthRatio": 0.7
  }
}
```

커스텀 config 파일은 필요한 부분만 작성하면 되며, 나머지는 기본값이 사용됩니다.

## 지원 색상 형식

- **Hex 코드**: `ff6b6b`, `#ff6b6b`
- **색상명**: `red`, `green`, `blue`, `yellow`, `purple`, `orange`, `pink`, `brown`, `gray`, `black`, `white`

## 지원 파일 형식

- SVG (`.svg`)
- PNG (`.png`)
- JPG/JPEG (`.jpg`, `.jpeg`)

## npm 스크립트

| 명령어 | 설명 |
|--------|------|
| `npm start` | 기본 config.json으로 썸네일 생성 |
| `npm run generate` | start와 동일 |
| `npm run bookreport` | config-bookreport.json으로 썸네일 생성 |
| `npm run dev` | watch 모드로 실행 (파일 변경 감지) |

새로운 config 파일을 위한 스크립트를 추가하려면 `package.json`의 `scripts` 섹션에 추가하세요:

```json
{
  "scripts": {
    "custom": "node generate.js -c ./config-custom.json"
  }
}
```

## 자동화 기능

1. **일괄 처리**: assets 폴더의 모든 로고 파일 자동 스캔
2. **파일명 파싱**: 로고명, 로고색상, 배경색상 자동 분석
3. **SVG 색상 변경**: SVG 파일의 fill/stroke 속성을 지정된 색상으로 자동 변경
4. **스마트 크기 조정**: maxWidth/maxHeight 기준으로 비율 유지하며 리사이즈
5. **중앙 배치**: 로고를 이미지 중앙에 자동 배치
6. **타겟 복사**: 설정된 경로로 자동 복사 (선택사항)
