package com.moseoh.searchencrypteddata;

import jakarta.persistence.EntityManager;
import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validator;
import jakarta.validation.constraints.Size;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.transaction.annotation.Transactional;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import com.moseoh.searchencrypteddata.entity.TestEntity;
import com.moseoh.searchencrypteddata.entity.TestEntityRepository;
import com.moseoh.searchencrypteddata.utils.HashUtils;
import java.util.List;
import java.util.Optional;
import java.util.Set;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
class TestEntityRepositoryTest {

    @Autowired
    private TestEntityRepository repository;

    @Autowired
    private EntityManager em;

    @Autowired
    private Validator validator;

    @AfterEach
    void tearDown() {
        repository.deleteAll();
    }

    @Test
    @DisplayName("글자 수가 초과하여 Exception이 발생한다.")
    void 엔티티_생성_테스트() {
        // given - fullName이 3글자를 초과하는 케이스 (홍길동동 = 4글자)
        TestEntity newEntity = new TestEntity("길동동", "홍");

        // when - validator로 직접 검증
        Set<ConstraintViolation<TestEntity>> violations = validator.validate(newEntity);

        // then - @Size(max = 3) 위반으로 인한 validation 오류가 발생해야 함
        assertThat(violations).isNotEmpty();
        List<String> violationNames = violations.stream().map(it -> it.getPropertyPath().toString()).toList();
        ConstraintViolation<TestEntity> violation = violations.iterator().next();
        assertThat(violation.getConstraintDescriptor().getAnnotation()).isInstanceOf(Size.class);
        assertThat(violationNames).contains("name", "fullName");
    }

    @Test
    @Transactional
    @DisplayName("같은 값을 저장해도 암호화된 값은 다르다.")
    void 엔티티_생성_테스트2() {
        TestEntity newEntity1 = new TestEntity("길동", "홍");
        TestEntity newEntity2 = new TestEntity("길동", "홍");

        repository.save(newEntity1);
        repository.save(newEntity2);
        em.flush();
        em.clear();

        // then - Native Query를 사용해 DB의 암호화된 값을 직접 조회
        // 참고: 테이블명(test_entity)과 컬럼명(full_name)은 실제 DB 생성 규칙에 맞게 조정하세요.
        String encryptedFullName1 = (String) em
                .createNativeQuery("SELECT full_name FROM test_entity WHERE id = :id")
                .setParameter("id", newEntity1.getId())
                .getSingleResult();

        String encryptedFullName2 = (String) em
                .createNativeQuery("SELECT full_name FROM test_entity WHERE id = :id")
                .setParameter("id", newEntity2.getId())
                .getSingleResult();

        // 검증 1: 두 암호화된 값은 null이 아니어야 한다.
        assertThat(encryptedFullName1).isNotNull();
        assertThat(encryptedFullName2).isNotNull();

        // 검증 2: 두 암호화된 값은 서로 달라야 한다. (핵심)
        assertThat(encryptedFullName1).isNotEqualTo(encryptedFullName2);

        // 검증 3 : 복호화된 두 값은 같아야 한다.
        TestEntity fetchedEntity1 = repository.findById(newEntity1.getId()).orElseThrow();
        TestEntity fetchedEntity2 = repository.findById(newEntity2.getId()).orElseThrow();
        assertThat(fetchedEntity1.getFullName()).isEqualTo(fetchedEntity2.getFullName());
    }

    @Test
    @DisplayName("저장시 hash 값이 생성되는지 확인")
    void 엔티티_생성_테스트3() {
        TestEntity newEntity = new TestEntity("길동", "홍");

        TestEntity savedEntity = repository.save(newEntity);

        assertThat(savedEntity.getNameHash()).isNotNull();
        assertThat(savedEntity.getFamilyNameHash()).isNotNull();
        assertThat(savedEntity.getFullNameHash()).isNotNull();
    }

    @Test
    @DisplayName("동일한 값들은 hash 값도 같은지 확인")
    void 엔티티_생성_테스트4() {
        TestEntity newEntity1 = new TestEntity("길동", "홍");
        TestEntity newEntity2 = new TestEntity("길동", "홍");

        TestEntity savedEntity1 = repository.save(newEntity1);
        TestEntity savedEntity2 = repository.save(newEntity2);

        assertThat(savedEntity1.getNameHash()).isEqualTo(savedEntity2.getNameHash());
        assertThat(savedEntity1.getFamilyNameHash()).isEqualTo(savedEntity2.getFamilyNameHash());
        assertThat(savedEntity1.getFullNameHash()).isEqualTo(savedEntity2.getFullNameHash());
    }

    @Test
    @DisplayName("hash 값으로 검색시 올바르게 반환한다.")
    void 이름_해시로_조회_테스트() {
        // given
        HashUtils hashUtils = new HashUtils();
        TestEntity newEntity1 = new TestEntity("길동", "홍");
        repository.save(newEntity1);
        String searchName = "길동";

        // when
        String nameHash = hashUtils.generate(searchName);
        Optional<TestEntity> foundEntity = repository.findByNameHash(nameHash);

        // then
        assertThat(foundEntity).isPresent();
        assertThat(foundEntity.get().getName()).isEqualTo(searchName);
        assertThat(foundEntity.get().getNameHash()).isEqualTo(nameHash);
    }

}