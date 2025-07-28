package com.moseoh.searchencrypteddata.entity;

import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import com.moseoh.searchencrypteddata.utils.HashUtils;

@Component
@RequiredArgsConstructor
public class TestEntityListener {

    private final HashUtils hashUtils;

    @PrePersist
    @PreUpdate
    public void generateHashes(TestEntity entity) {
        if (StringUtils.hasText(entity.getName())) {
            entity.setNameHash(this.hashUtils.generate(entity.getName()));
        }

        if (StringUtils.hasText(entity.getFamilyName())) {
            entity.setFamilyNameHash(this.hashUtils.generate(entity.getFamilyName()));
        }

        if (StringUtils.hasText(entity.getFullName())) {
            entity.setFullNameHash(this.hashUtils.generate(entity.getFullName()));
        }
    }

}
