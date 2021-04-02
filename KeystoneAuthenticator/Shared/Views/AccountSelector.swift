//
//  AccountSelector.swift
//  KeystoneAuthenticator
//
//  Created by Aaron Craelius on 4/1/21.
//

import SwiftUI

struct AccountSelector: View {
    @FetchRequest(entity: Account.entity(), sortDescriptors: []) var accounts: FetchedResults<Account>
    
    let addAccount: () -> Void
    let onSelected: (Account) -> Void
    
    @State var inAddAccountView = true
    
    var body: some View {
        NavigationView {
            VStack {
                Text("Select An Account").font(/*@START_MENU_TOKEN@*/.title/*@END_MENU_TOKEN@*/)
                
                ForEach(accounts, id: \.id) { account in
                    Button(action: {onSelected(account)}) {
                        Text(account.name ?? "Unknown")
                    }
                }.padding(5)
                .border(Color.black, width: /*@START_MENU_TOKEN@*/1/*@END_MENU_TOKEN@*/)
                
                NavigationLink(destination:WelcomeView(onDone: {})) {
                    Label("Add Account", systemImage: "plus")
                }.padding(5)
            }
        }
    }
}

struct AccountSelector_Previews: PreviewProvider {
    static var previews: some View {
        AccountSelector(addAccount: {}, onSelected: { account in })
            .environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
