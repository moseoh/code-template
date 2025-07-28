package com.moseoh.searchencrypteddata.utils;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

import org.springframework.stereotype.Component;

import com.moseoh.searchencrypteddata.SecurityProperties;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

/**
 * Hash 생성시 항상 동일한 길이를 갖습니다. => 44자리
 */
@Component
public class HashUtils {
    private static final String        HMAC_ALGORITHM = "HmacSHA256";
    private final        SecretKeySpec hmacKey;

    // application.yml에 정의된 별도의 hmac 키를 주입받습니다.
    public HashUtils() {
        byte[] secretKeyBytes = Base64.getDecoder().decode(SecurityProperties.SECRET_KEY);
        this.hmacKey = new SecretKeySpec(secretKeyBytes, HMAC_ALGORITHM);
    }

    public String generate(String data) {
        if (data == null) {
            return null;
        }
        try {
            Mac mac = Mac.getInstance(HMAC_ALGORITHM);
            mac.init(hmacKey);
            byte[] hash = mac.doFinal(data.getBytes());
            return Base64.getEncoder().encodeToString(hash);
        } catch (NoSuchAlgorithmException |
                 InvalidKeyException e) {
            // 실제 프로덕션에서는 복구 불가능한 심각한 오류로 처리해야 합니다.
            throw new IllegalStateException("Failed to generate HMAC.", e);
        }
    }
}
