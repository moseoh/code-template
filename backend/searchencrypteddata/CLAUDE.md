# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 프로젝트 개요

이 프로젝트는 Spring Boot 기반의 암호화된 데이터 검색 시스템입니다. AES-GCM 암호화를 사용하여 민감한 데이터를 안전하게 저장하고 해시를 통한 검색 기능을 제공합니다.

## 개발 명령어

### 프로젝트 빌드 및 실행
```bash
# 프로젝트 빌드
./gradlew build

# 애플리케이션 실행
./gradlew bootRun

# 테스트 실행
./gradlew test

# 단일 테스트 실행
./gradlew test --tests "클래스명.메소드명"
```

### 데이터베이스 설정
```bash
# MariaDB 컨테이너 시작
docker-compose up -d

# 컨테이너 중지
docker-compose down
```

### 기타 유용한 명령어
```bash
# 의존성 확인
./gradlew dependencies

# 프로젝트 정리
./gradlew clean
```

## 아키텍처 및 주요 컴포넌트

### 핵심 기술 스택
- **Framework**: Spring Boot 3.5.3
- **Database**: MariaDB 11.6 (Docker Compose로 관리)
- **ORM**: Spring Data JPA + QueryDSL
- **암호화**: AES-GCM (256-bit)
- **Build Tool**: Gradle with Kotlin DSL

### 주요 컴포넌트

#### 1. AESConverter (src/main/java/com/moseoh/searchencrypteddata/AESConverter.java)
- JPA AttributeConverter 구현
- AES-GCM 암호화/복호화 담당
- 데이터베이스 저장 시 자동 암호화, 조회 시 자동 복호화
- 매번 새로운 IV(초기화 벡터) 생성으로 보안성 강화

#### 2. SecurityProperties (src/main/java/com/moseoh/searchencrypteddata/SecurityProperties.java)
- 암호화 키 관리 클래스
- Base64로 인코딩된 256-bit 암호화 키 저장
- **보안 중요**: 실제 운영 환경에서는 환경변수나 별도 키 관리 시스템 사용 필요

#### 3. TestEntity (src/main/java/com/moseoh/searchencrypteddata/TestEntity.java)
- 암호화된 검색 기능 구현 예시 엔티티
- name, familyName, fullName 필드는 암호화되어 저장
- nameHash, familyNameHash, fullNameHash 필드는 검색용 해시값 저장
- @PrePersist/@PreUpdate 훅을 통한 해시값 자동 업데이트

### 데이터베이스 설정
- **연결**: localhost:63306 (Docker 컨테이너)
- **데이터베이스명**: search-encrypted-data
- **사용자**: user / password
- **Hibernate DDL**: create-drop (개발용)

### 보안 고려사항
1. **암호화 키 관리**: 현재 하드코딩된 키는 개발용이며, 운영 환경에서는 보안 강화 필요
2. **해시 기반 검색**: 암호화된 필드에 대한 검색을 위해 해시값 활용
3. **AES-GCM 모드**: 인증과 암호화를 동시에 제공하는 강력한 암호화 방식 사용

### QueryDSL 설정
- QueryDSL JPA 5.1.0 설정 완료
- 복잡한 쿼리 작성 시 타입 안전성 보장
- Blaze Persistence 통합으로 고급 쿼리 기능 지원

## 개발 시 주의사항

1. **암호화 필드 추가 시**: AESConverter를 @Convert 어노테이션으로 적용
2. **검색 기능 구현 시**: 암호화된 필드는 직접 검색 불가, 해시 필드 활용
3. **테스트 작성 시**: 암호화/복호화 로직 검증 필수
4. **데이터베이스 스키마 변경 시**: DDL 모드가 create-drop이므로 데이터 손실 주의