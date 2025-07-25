package com.moseoh.searchencrypteddata;

import jakarta.validation.ConstraintViolation;
import jakarta.validation.Validator;
import jakarta.validation.constraints.Size;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import org.junit.jupiter.api.Test;
import java.util.Set;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
class TestEntityRepositoryTest {

    @Autowired
    private TestEntityRepository testEntityRepository;

    @Autowired
    private Validator validator;

    @Test
    void 엔티티_생성_테스트() {
        // given - fullName이 3글자를 초과하는 케이스 (박민수2 = 4글자)
        TestEntity newEntity = new TestEntity("민수2", "박");

        // when - validator로 직접 검증
        Set<ConstraintViolation<TestEntity>> violations = validator.validate(newEntity);

        // then - @Size(max = 3) 위반으로 인한 validation 오류가 발생해야 함
        assertThat(violations).isNotEmpty();
        ConstraintViolation<TestEntity> violation = violations.iterator().next();
        assertThat(violation.getConstraintDescriptor().getAnnotation()).isInstanceOf(Size.class);
        assertThat(violation.getPropertyPath().toString()).isEqualTo("fullName");
    }

    //    @Test
    //    void 이름_해시로_조회_테스트() {
    //        // given
    //        String searchNameHash = "홍길동";
    //
    //        // when
    //        Optional<TestEntity> foundEntity = testEntityRepository.findByNameHash(searchNameHash);
    //
    //        // then
    //        assertThat(foundEntity).isPresent();
    //        assertThat(foundEntity.get().getName()).isEqualTo("홍길동");
    //        assertThat(foundEntity.get().getNameHash()).isEqualTo(searchNameHash);
    //    }
    //
    //    @Test
    //    void 성_해시로_조회_테스트() {
    //        // given
    //        String searchFamilyNameHash = "김";
    //
    //        // when
    //        Optional<TestEntity> foundEntity = testEntityRepository.findByFamilyNameHash(searchFamilyNameHash);
    //
    //        // then
    //        assertThat(foundEntity).isPresent();
    //        assertThat(foundEntity.get().getFamilyName()).isEqualTo("김");
    //        assertThat(foundEntity.get().getFamilyNameHash()).isEqualTo(searchFamilyNameHash);
    //    }
    //
    //    @Test
    //    void 전체이름_해시로_조회_테스트() {
    //        // given
    //        String searchFullNameHash = "이영희이";
    //
    //        // when
    //        Optional<TestEntity> foundEntity = testEntityRepository.findByFullNameHash(searchFullNameHash);
    //
    //        // then
    //        assertThat(foundEntity).isPresent();
    //        assertThat(foundEntity.get().getFullName()).isEqualTo("이영희이");
    //        assertThat(foundEntity.get().getFullNameHash()).isEqualTo(searchFullNameHash);
    //    }
    //
    //    @Test
    //    void 이름과_성_해시로_조회_테스트() {
    //        // given
    //        String nameHash = "홍길동";
    //        String familyNameHash = "홍";
    //
    //        // when
    //        List<TestEntity> foundEntities = testEntityRepository.findByNameHashAndFamilyNameHash(nameHash, familyNameHash);
    //
    //        // then
    //        assertThat(foundEntities).hasSize(1);
    //        assertThat(foundEntities.get(0).getName()).isEqualTo("홍길동");
    //        assertThat(foundEntities.get(0).getFamilyName()).isEqualTo("홍");
    //    }
    //
    //    @Test
    //    void 모든_엔티티_조회_테스트() {
    //        // when
    //        List<TestEntity> allEntities = testEntityRepository.findAll();
    //
    //        // then
    //        assertThat(allEntities).hasSize(3);
    //        assertThat(allEntities).extracting(TestEntity::getName).containsExactlyInAnyOrder("홍길동", "김철수", "이영희");
    //    }
    //
    //    @Test
    //    void 엔티티_삭제_테스트() {
    //        // given
    //        Long initialCount = testEntityRepository.count();
    //
    //        // when
    //        testEntityRepository.delete(testEntity1);
    //
    //        // then
    //        assertThat(testEntityRepository.count()).isEqualTo(initialCount - 1);
    //        assertThat(testEntityRepository.findByNameHash("홍길동")).isEmpty();
    //    }
    //
    //    @Test
    //    void 부분_검색_테스트() {
    //        // given - 이름에 "길" 이 포함된 엔티티 검색
    //        String partialNameHash = "길";
    //
    //        // when
    //        List<TestEntity> foundEntities = testEntityRepository.findByNameHashContaining(partialNameHash);
    //
    //        // then
    //        assertThat(foundEntities).hasSize(1);
    //        assertThat(foundEntities.get(0).getName()).isEqualTo("홍길동");
    //    }
}