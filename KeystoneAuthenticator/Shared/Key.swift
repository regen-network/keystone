//
//  Key.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import Foundation
import CryptoKit
import Security
import CoreData


let CurrentAccountKey = "currentAccount"

func getCurrentAccount() -> String? {
    UserDefaults.standard.string(forKey: CurrentAccountKey)
}

func setCurrentAccount(account: Account) {
    UserDefaults.standard.setValue(account.objectID.uriRepresentation().absoluteString, forKeyPath: CurrentAccountKey)
}

enum AccountError : Error {
    case keychainAdd
    case save
}

func createAccount(name: String, chainId: String, moc: NSManagedObjectContext) throws {
    // create new private key
    //let privKey = try SecureEnclave.P256.Signing.PrivateKey()
    // store priv key in keychain under account name
//    let query: [String: Any] = [kSecClass as String: kSecClassKey,
//                                kSecAttrLabel as String: name,
//                                kSecAttrApplicationLabel as String: chainId,
//                                kSecValueData as String: privKey.dataRepresentation]
//
//    let status = SecItemAdd(query as CFDictionary, nil)
//    guard status == errSecSuccess else { throw AccountError.keychainAdd }
    
    let account = Account(context: moc)
    account.name = name
    account.chainId = chainId
    
    do {
        try moc.save()
    } catch {
        throw AccountError.save
    }
    
    setCurrentAccount(account: account)
}

func getAccounts() {
    let query: [String: Any] = [kSecClass as String: kSecClassKey,
                                kSecMatchLimit as String: kSecMatchLimitAll,
                                kSecReturnAttributes as String: true,
                                kSecReturnData as String: true]
    
    var item: CFTypeRef?
    let status = SecItemCopyMatching(query as CFDictionary, &item)
//    guard status != errSecItemNotFound else { throw KeychainError.noPassword }
//    guard status == errSecSuccess else { throw KeychainError.unhandledError(status: status) }
}
