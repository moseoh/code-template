# 개인정보 암호화 및 해시 검색 프로젝트

## 핵심 기능

### 개인정보 암호화

- **JPA AttributeConverter 활용:** `AESConverter` 클래스는 JPA의 `AttributeConverter` 인터페이스를 구현하여, 엔티티의 특정 필드(예: `name`, `familyName`, `fullName`)가 데이터베이스에 저장되거나 조회될 때 자동으로 암호화 및 복호화를 수행합니다.
- **AES-256 암호화:** `AESUtils` 클래스를 통해 AES-256 알고리즘을 사용하여 데이터를 암호화합니다. 암호화 키는 `SecurityProperties`에 정의된 값을 사용하며, 매번 다른 IV(Initialization Vector)를 사용하여 동일한 데이터라도 암호화된 결과가 다르게 나옵니다. 이는 보안을 강화하는 중요한 요소입니다.

### 해시 기반 검색

- **HMAC-SHA256 해시 생성:** `HashUtils` 클래스는 HMAC-SHA256 알고리즘을 사용하여 데이터의 해시 값을 생성합니다. 해시 키는 `SecurityProperties`에 정의된 별도의 키를 사용합니다. 동일한 데이터는 항상 동일한 해시 값을 생성하므로, 이를 통해 암호화된 데이터를 검색할 수 있습니다.
- **JPA EntityListener 활용:** `TestEntityListener`는 `@PrePersist`와 `@PreUpdate` 어노테이션을 사용하여 엔티티가 저장되거나 업데이트되기 전에 자동으로 해시 값을 생성하여 해당 필드에 저장합니다.
- **해시 값으로 검색:** `TestEntityRepository`에서 `findByNameHash`와 같은 메서드를 통해 해시 값을 사용하여 데이터를 조회합니다. 이를 통해 원본 데이터를 노출하지 않고도 안전하게 검색할 수 있습니다.

## 프로젝트 구조

- **`com.moseoh.searchencrypteddata`**: 메인 패키지
  - **`entity`**: 데이터베이스 테이블과 매핑되는 엔티티 클래스
    - `TestEntity`: 이름, 성, 전체 이름을 저장하는 엔티티. `@Convert` 어노테이션을 통해 `AESConverter`를 사용하고, `@EntityListeners`를 통해 `TestEntityListener`를 사용합니다.
    - `TestEntityListener`: 엔티티 저장/수정 시 해시 값을 생성하는 리스너.
    - `TestEntityRepository`: `TestEntity`에 대한 CRUD 및 해시 검색을 위한 JpaRepository.
  - **`utils`**: 암호화 및 해시 관련 유틸리티 클래스
    - `AESUtils`: AES-256 암호화 및 복호화 유틸리티.
    - `HashUtils`: HMAC-SHA256 해시 생성 유틸리티.
  - **`AESConverter`**: JPA AttributeConverter 구현체.
  - **`SecurityProperties`**: 암호화 및 해시 키를 관리하는 클래스.

## 테스트

- **`TestEntityRepositoryTest`**: `TestEntity`의 저장, 조회, 암호화, 해시 생성 및 검색 기능에 대한 단위 테스트를 포함합니다.
  - **글자 수 초과 테스트:** `@Size` 어노테이션을 이용한 유효성 검증 테스트.
  - **암호화 값 비교 테스트:** 동일한 값을 저장해도 암호화된 값은 다른지 확인.
  - **해시 값 생성 테스트:** 저장 시 해시 값이 올바르게 생성되는지 확인.
  - **해시 값 일관성 테스트:** 동일한 값은 동일한 해시 값을 갖는지 확인.
  - **해시 값 검색 테스트:** 해시 값으로 데이터를 올바르게 조회하는지 확인.

- **`ValidationSchemaTest`**: AES 암호화 후 데이터 길이를 검증합니다. `@Size` 어노테이션의 `max` 값에 따라 한글을 기준으로 암호화된 결과물의 길이가 `@Column`의 `length`를 초과하는지 확인하여, 데이터베이스 스키마의 유효성을 보장합니다.

## 중요 설정

- **`spring.jpa.properties.hibernate.validator.apply_to_ddl=false`**: `src/test/resources/application.properties` 파일의 이 설정은 Hibernate가 DDL(Data Definition Language)을 생성할 때 `@Size`와 같은 Bean Validation 어노테이션을 무시하도록 합니다. 이를 통해 암호화된 데이터의 길이를 기준으로 `@Column(length=...)`를 설정할 수 있으며, 원본 데이터의 길이를 기준으로 하는 `@Size` 어노테이션에 의해 DB 컬럼 길이가 제한되는 문제를 방지합니다. 자세한 설명은 해당 파일에 적어두었습니다.