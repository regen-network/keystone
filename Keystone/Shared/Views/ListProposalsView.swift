//
//  ListProposalsView.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import SwiftUI

let exampleTxs =  [
    TxInfo(summary: "Send 100REGEN from regen1abcdecxvad49sdf2 to regen1sdgkhwfeiaflknewio"),
    TxInfo(summary: "Sign regen:HsjWGj592kSw924FDWKeIWH"),
]

struct ListProposalsView: View {
    

    var txs: [TxInfo] = exampleTxs

    @State private var multiSelection = Set<UUID>()

    var body: some View {
        NavigationView {
            List {
                ForEach(txs) { tx in
                    NavigationLink(destination: ApproveTxView(tx: tx)) {
                        Text(tx.summary)
                    }
                }
            }.navigationBarItems(trailing: NavigationLink(destination: SettingsView()) {
                Image(systemName: "gear")
            })
        }
    }
}

struct ListProposalsView_Previews: PreviewProvider {
        
    static var previews: some View {
        ListProposalsView(txs: exampleTxs)
    }
}
