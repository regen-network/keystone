//
//  WelcomeView.swift
//  KeystoneAuthenticator
//
//  Created by Aaron Craelius on 3/31/21.
//

import SwiftUI


enum WelcomeState {
    case creating
    case done
}


struct WelcomeView: View {
    @State var selectedChain = ""
    @State var accountName = ""

    @State var isCreating = false
    @State var isDone = false
    @State var isError = false
    
    @Environment(\.managedObjectContext) var moc
    
    let onDone: () -> Void

    var body: some View {
        VStack {
            NavigationLink(destination: VStack {
                Picker("Select Network", selection: $selectedChain) {
                    Text("Cosmos Hub")
                    Text("Regen Network")
                }.navigationTitle("Select Network")

                NavigationLink(destination: VStack {
                    TextField("Account Name", text: $accountName).padding(20)
                            .multilineTextAlignment(.center)
                    Button(action: {
                        do {
                            isCreating = true
                            try createAccount(name: accountName, chainId: selectedChain, moc: self.moc)
                            isDone = true
                            onDone()
                        } catch {
                            isError = true
                        }
                    }) {
                        Text("Create Account")
                    }.disabled(accountName == "")
                }) {
                    HStack {
                        Text("Next")
                    }
                }
            }) {
                Text("Create New Account")
            }.padding(10)
            Button(action: {}) {
                Text("Sign-in to Existing Account")
            }.padding(10)
        }
                .navigationTitle("Add Account")
                .fullScreenCover(isPresented: $isCreating, content: {
                    VStack {
                        Text("Creating Account...")
                        ProgressView().progressViewStyle(CircularProgressViewStyle())
                    }
                })
                .fullScreenCover(isPresented: $isDone, content: {
                    VStack {
                        Text("Done")
                    }                })
                .alert(isPresented: $isError, content: {
                    Alert(title: Text("Error"))
                })
    }
}

struct WelcomeView_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView {
            WelcomeView(onDone: {})
        }
    }
}
