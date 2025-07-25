package com.moseoh.searchencrypteddata;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface TestEntityRepository extends JpaRepository<TestEntity, UUID> {

    Optional<TestEntity> findByNameHash(String nameHash);

    Optional<TestEntity> findByFamilyNameHash(String familyNameHash);

    Optional<TestEntity> findByFullNameHash(String fullNameHash);

    List<TestEntity> findByNameHashContaining(String nameHash);

    List<TestEntity> findByFamilyNameHashContaining(String familyNameHash);

    @Query("SELECT t FROM TestEntity t WHERE t.nameHash = :nameHash AND t.familyNameHash = :familyNameHash")
    List<TestEntity> findByNameHashAndFamilyNameHash(@Param("nameHash") String nameHash, @Param("familyNameHash") String familyNameHash);
}