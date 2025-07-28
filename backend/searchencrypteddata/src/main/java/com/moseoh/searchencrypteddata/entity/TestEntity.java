package com.moseoh.searchencrypteddata.entity;

import jakarta.persistence.*;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.NoArgsConstructor;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.moseoh.searchencrypteddata.AESConverter;
import java.util.UUID;

/**
 * AES 암호화시 길이에 따라, 한글 & 영어에 따라 다른 길이를 갖습니다.
 * {@link com.moseoh.searchencrypteddata.ValidationSchemaTest} 을 통해 DB Column 길이를 확인합니다.
 * <p>
 * Hash 생성시 항상 동일한 길이를 갖습니다. => 44자리
 */
@Entity
@Table(name = "test_entity")
@Getter
@NoArgsConstructor
@EntityListeners(TestEntityListener.class)
public class TestEntity {

    @Id
    private UUID id;

    @Size(max = 2)
    @Column(name = "name", length = 48, nullable = false)
    @Convert(converter = AESConverter.class)
    private String name;

    @Size(max = 1)
    @Column(name = "family_name", length = 44, nullable = false)
    @Convert(converter = AESConverter.class)
    private String familyName;

    @Size(max = 3)
    @Column(name = "full_name", length = 52, nullable = false)
    @Convert(converter = AESConverter.class)
    private String fullName;

    @Column(name = "name_hash", length = 44, nullable = false)
    private String nameHash;

    @Column(name = "family_name_hash", length = 44, nullable = false)
    private String familyNameHash;

    @Column(name = "full_name_hash", length = 44, nullable = false)
    private String fullNameHash;

    public TestEntity(
            String name,
            String familyName
    ) {
        this.id = UUID.randomUUID();
        this.name = name;
        this.familyName = familyName;
        this.fullName = familyName + name;
    }

    @Override
    public String toString() {
        ObjectMapper mapper = new ObjectMapper();
        mapper.enable(SerializationFeature.INDENT_OUTPUT);
        try {
            return mapper.writeValueAsString(this);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    protected void setNameHash(String nameHash) {
        this.nameHash = nameHash;
    }

    protected void setFamilyNameHash(String familyNameHash) {
        this.familyNameHash = familyNameHash;
    }

    protected void setFullNameHash(String fullNameHash) {
        this.fullNameHash = fullNameHash;
    }
}
