package com.moseoh.searchencrypteddata;

import jakarta.persistence.*;
import jakarta.validation.constraints.Size;
import lombok.Getter;
import lombok.NoArgsConstructor;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import java.util.UUID;

@Entity
@Table(name = "test_entity")
@Getter
@NoArgsConstructor
public class TestEntity {

    @Id
    private UUID id;

    @Column(name = "name", length = 255, nullable = false)
    private String name;

    @Column(name = "family_name", length = 255, nullable = false)
    private String familyName;

    @Size(max = 53)
    @Column(name = "full_name", length = 255, nullable = false)
    @Convert(converter = AESConverter.class)
    private String fullName;

    @Column(name = "name_hash", length = 255, nullable = false)
    private String nameHash;

    @Column(name = "family_name_hash", length = 255, nullable = false)
    private String familyNameHash;

    @Column(name = "full_name_hash", length = 255, nullable = false)
    private String fullNameHash;

    @PrePersist
    @PreUpdate
    private void updateHash() {
        this.nameHash = this.name;
        this.familyNameHash = this.familyName;
        this.fullNameHash = this.fullName;
    }

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
}
