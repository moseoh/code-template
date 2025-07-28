package com.moseoh.searchencrypteddata;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;
import lombok.RequiredArgsConstructor;

import com.moseoh.searchencrypteddata.utils.AESUtils;

/**
 * AES 암호화시 길이에 따라, 한글 & 영어에 따라 다른 길이를 갖습니다.
 * {@link com.moseoh.searchencrypteddata.ValidationSchemaTest} 을 통해 DB Column 길이를 확인합니다.
 */
@Converter
@RequiredArgsConstructor
public class AESConverter implements AttributeConverter<String, String> {

    private final AESUtils aesUtils;

    @Override
    public String convertToDatabaseColumn(String attribute) {
        if (attribute == null || attribute.isEmpty()) {
            return attribute;
        }
        return aesUtils.encrypt(attribute);
    }

    @Override
    public String convertToEntityAttribute(String dbData) {
        if (dbData == null || dbData.isEmpty()) {
            return dbData;
        }
        return aesUtils.decrypt(dbData);
    }

}