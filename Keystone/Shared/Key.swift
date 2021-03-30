//
//  Key.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import Foundation

func createKey() throws -> SecKey {
    let access =
        SecAccessControlCreateWithFlags(
            kCFAllocatorDefault,
            kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly,
            .privateKeyUsage,
            nil)
    
    let tag = "network.regen.keystone.Key1".data(using: .utf8)!
    
    let attributes: [String: Any] = [
      kSecAttrKeyType as String:            kSecAttrKeyTypeECSECPrimeRandom,
      kSecAttrKeySizeInBits as String:      256,
      kSecAttrTokenID as String:            kSecAttrTokenIDSecureEnclave,
      kSecPrivateKeyAttrs as String: [
        kSecAttrIsPermanent as String:      true,
        kSecAttrApplicationTag as String:   tag,
        kSecAttrAccessControl as String:    access
      ]
    ]
    
    var error: Unmanaged<CFError>?
    guard let privateKey = SecKeyCreateRandomKey(attributes as CFDictionary, &error) else {
        throw error!.takeRetainedValue() as Error
    }

    return privateKey
}

func getKey() -> SecKey? {
    let tag = "network.regen.keystone.Key1".data(using: .utf8)!


    let getquery: [String: Any] = [kSecClass as String: kSecClassKey,
                                   kSecAttrApplicationTag as String: tag,
                                   kSecAttrKeyType as String: kSecAttrKeyTypeECSECPrimeRandom,
                                   kSecReturnRef as String: true]
    
    var item: CFTypeRef?
    let status = SecItemCopyMatching(getquery as CFDictionary, &item)
    guard status == errSecSuccess else { return nil }
    return item as! SecKey
}

func getOrCreateKey() throws -> SecKey {
    try getKey() ?? createKey()
}
