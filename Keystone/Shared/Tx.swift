//
//  Tx.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import Foundation

struct TxInfo: Identifiable, Hashable {
    let summary: String
    let id = UUID()
}
