package com.moseoh.searchencrypteddata;

import jakarta.persistence.Column;
import jakarta.persistence.Convert;
import jakarta.persistence.Entity;
import jakarta.validation.constraints.Size;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;
import org.reflections.Reflections;

import java.lang.reflect.Field;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.assertTrue;

class ValidationSchemaTest {

    private static final String ENTITY_PACKAGE = "com.moseoh.searchencrypteddata";

    static Stream<Class<?>> entityProvider() {
        Reflections reflections = new Reflections(ENTITY_PACKAGE);
        return reflections.getTypesAnnotatedWith(Entity.class).stream();
    }

    /**
     * 엔티티의 암호화 필드 스키마를 검증합니다.
     * <p>
     * 지정된 패키지 내의 모든 {@code @Entity}를 스캔하여,
     * {@code @Size}, {@code @Column}, {@code @Convert(converter = AESConverter.class)} 어노테이션을
     * 모두 가진 필드를 찾아냅니다.
     * <p>
     * {@code @Size}는 글자 수를 기준으로 검사하지만, 암호화 후의 길이는 바이트에 따라 달라집니다.
     * 한글은 UTF-8 인코딩 시 영문/숫자보다 더 많은 바이트를 차지하므로, 암호화 후 결과 문자열이 더 길어집니다.
     * 따라서 가장 긴 결과가 나오는 <b>한글('기')을 기준</b>으로 테스트하여 DB 컬럼의 길이가 충분한지 엄격하게 확인합니다.
     *
     * @param entityClass 검증할 엔티티 클래스
     * @throws Exception 리플렉션 또는 컨버터 생성 시 발생할 수 있는 예외
     */
    @DisplayName("암호화 필드 검증: AESConverter의 @Size 대비 @Column(length)가 충분한가")
    @ParameterizedTest(name = "{0} 엔티티 검증")
    @MethodSource("entityProvider")
    void validateEncryptedColumnLength(Class<?> entityClass) throws Exception {
        for (Field field : entityClass.getDeclaredFields()) {
            Size sizeAnn = field.getAnnotation(Size.class);
            Column columnAnn = field.getAnnotation(Column.class);
            Convert convertAnn = field.getAnnotation(Convert.class);

            if (sizeAnn == null || columnAnn == null || convertAnn == null) {
                continue;
            }

            if (!convertAnn.converter().equals(AESConverter.class)) {
                continue;
            }

            AESConverter converter = new AESConverter();

            String testInput = "기".repeat(sizeAnn.max());
            String encryptedOutput = converter.convertToDatabaseColumn(testInput);

            String errorMessage = String.format(
                    "검증 실패! [%s.%s]: 원본 최대 길이(%d) -> 암호화된 길이(%d) > DB 컬럼 길이(%d)",
                    entityClass.getSimpleName(), field.getName(),
                    sizeAnn.max(), encryptedOutput.length(), columnAnn.length()
            );

            assertTrue(encryptedOutput.length() <= columnAnn.length(), errorMessage);

            System.out.printf(
                    "검증 통과 ✅: [%s.%s] (원본: %d, 암호화: %d, DB컬럼: %d)%n",
                    entityClass.getSimpleName(),
                    field.getName(),
                    sizeAnn.max(),
                    encryptedOutput.length(),
                    columnAnn.length()
            );
        }
    }
}