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

    @AppStorage(CurrentAccountKey)
    private var selectedAcct: String = ""

    @State private var navigateTo = ""

    @State private var inSettingsView = false
    
    @State private var inAddAccountView = false
    
    @FetchRequest(entity: Account.entity(), sortDescriptors: []) var accounts: FetchedResults<Account>
    
    func selectedAccount() -> Account? {
        for account in accounts {
            if account.objectID.uriRepresentation().absoluteString == selectedAcct {
                return account
            }
        }
        
        return nil
    }
        
    var body: some View {
        if accounts.count == 0 {
            NavigationView {
                WelcomeView(onDone: {})
            }
        } else if selectedAccount() == nil {
            AccountSelector(addAccount: {
                
            }, onSelected: {
                account in selectedAcct = account.objectID.uriRepresentation().absoluteString
            })
        } else {
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
                                            ForEach(accounts, id: \.id) { account in
                                                Text(account.name ?? "Unknown").tag(account.objectID.uriRepresentation().absoluteString)
                                            }
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
                                WelcomeView(onDone: {
                                    inAddAccountView = false
                                    
                                }).toolbar {
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
}

struct ListProposalsView_Previews: PreviewProvider {

    static var previews: some View {
        ListTxsView(txs: exampleTxs)
            .environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
