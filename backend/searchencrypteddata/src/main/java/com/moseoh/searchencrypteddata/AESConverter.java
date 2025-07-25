package com.moseoh.searchencrypteddata;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.Base64;
import javax.crypto.Cipher;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class AESConverter implements AttributeConverter<String, String> {

    private static final int    GCM_IV_LENGTH            = 12;
    private static final int    GCM_TAG_LENGTH           = 128;
    private static final String ALGORITHM_TRANSFORMATION = "AES/GCM/NoPadding";
    private static final String ALGORITHM                = "AES";

    private final SecretKeySpec secretKeySpec;

    public AESConverter() {
        byte[] secretKeyBytes = Base64.getDecoder().decode(SecurityProperties.SECRET_KEY);
        // SecretKeySpec 생성 시에는 순수 알고리즘 이름("AES")을 사용합니다.
        this.secretKeySpec = new SecretKeySpec(secretKeyBytes, ALGORITHM);
    }

    @Override
    public String convertToDatabaseColumn(String attribute) {
        if (attribute == null || attribute.isEmpty()) {
            return attribute;
        }
        return encrypt(attribute);
    }

    @Override
    public String convertToEntityAttribute(String dbData) {
        if (dbData == null || dbData.isEmpty()) {
            return dbData;
        }
        return decrypt(dbData);
    }

    private String encrypt(String value) {
        try {
            // 1. GCM 모드는 매 암호화마다 새로운 IV(초기화 벡터)를 생성해야 합니다.
            byte[] iv = new byte[GCM_IV_LENGTH];
            new SecureRandom().nextBytes(iv);

            Cipher cipher = Cipher.getInstance(ALGORITHM_TRANSFORMATION);
            GCMParameterSpec gcmParameterSpec = new GCMParameterSpec(GCM_TAG_LENGTH, iv);

            cipher.init(Cipher.ENCRYPT_MODE, secretKeySpec, gcmParameterSpec);
            byte[] encryptedBytes = cipher.doFinal(value.getBytes(StandardCharsets.UTF_8));

            // 2. 복호화를 위해 IV와 암호문을 합쳐서 저장합니다.
            ByteBuffer byteBuffer = ByteBuffer.allocate(iv.length + encryptedBytes.length);
            byteBuffer.put(iv);
            byteBuffer.put(encryptedBytes);

            return Base64.getEncoder().encodeToString(byteBuffer.array());
        } catch (Exception e) {
            throw new RuntimeException("Error encrypting value", e);
        }
    }

    private String decrypt(String value) {
        try {
            // 1. 저장된 Base64 문자열을 디코딩하여 (IV + 암호문) 형태의 바이트 배열로 변환합니다.
            byte[] decodedBytes = Base64.getDecoder().decode(value);
            ByteBuffer byteBuffer = ByteBuffer.wrap(decodedBytes);

            // 2. IV와 암호문을 분리합니다.
            byte[] iv = new byte[GCM_IV_LENGTH];
            byteBuffer.get(iv);
            byte[] encryptedBytes = new byte[byteBuffer.remaining()];
            byteBuffer.get(encryptedBytes);

            Cipher cipher = Cipher.getInstance(ALGORITHM_TRANSFORMATION);
            GCMParameterSpec gcmParameterSpec = new GCMParameterSpec(GCM_TAG_LENGTH, iv);

            cipher.init(Cipher.DECRYPT_MODE, secretKeySpec, gcmParameterSpec);
            byte[] decryptedBytes = cipher.doFinal(encryptedBytes);

            return new String(decryptedBytes, StandardCharsets.UTF_8);
        } catch (Exception e) {
            throw new RuntimeException("Error decrypting value", e);
        }
    }
}