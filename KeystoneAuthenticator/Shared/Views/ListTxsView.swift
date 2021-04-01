//
//  ListProposalsView.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import SwiftUI

let exampleTxs = [
    TxInfo(summary: "Send 100REGEN from regen1abcdecxvad49sdf2 to regen1sdgkhwfeiaflknewio"),
    TxInfo(summary: "Sign regen:HsjWGj592kSw924FDWKeIWH"),
]

struct ListTxsView: View {


    var txs: [TxInfo] = exampleTxs

    @State private var multiSelection = Set<UUID>()

    @State private var selectedAcct = "1"

    @State private var navigateTo = ""

    @State private var inSettingsView = false
    
    @State private var inAddAccountView = false

    var body: some View {
        NavigationView {
            List {
                ForEach(txs) { tx in
                    NavigationLink(destination: ApproveTxView(tx: tx)) {
                        Text(tx.summary)
                    }
                }
            }.toolbar {
                        ToolbarItem(placement: .primaryAction) {
                            Menu {
                                Section {
                                    Button(action: {
                                        inSettingsView = true
                                    }) {
                                        Label("Manage Account", systemImage: "gear")
                                    }
                                }
                                Section {
                                    Picker("Account", selection: $selectedAcct) {
                                        Label("Account 1", systemImage: "person").tag("1")
                                        Label("Account 2", systemImage: "person").tag("2")
                                    }
                                    
                                    Button(action: {
                                        inAddAccountView = true
                                    }) {
                                        Label("Add Account", systemImage: "plus")
                                    }
                                }
                            } label: {
                                Label("Account", systemImage: "person.circle")
                            }.background(
                                NavigationLink(destination: SettingsView(), isActive: $inSettingsView) {
                                    EmptyView()
                                }
                            )
                        }
                    }
                    .navigationTitle("Pending Transactions")
                    .fullScreenCover(isPresented: $inAddAccountView) {
                        NavigationView {
                            WelcomeView().toolbar {
                                ToolbarItem(placement: .primaryAction) {
                                    Button(action: {
                                        inAddAccountView = false
                                    }) {
                                        Text("Cancel").fontWeight(.semibold)
                                    }
                                }
                            }
                        }
                    }
        }
    }
}

struct ListProposalsView_Previews: PreviewProvider {

    static var previews: some View {
        ListTxsView(txs: exampleTxs)
    }
}
