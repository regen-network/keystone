//
//  WelcomeView.swift
//  KeystoneAuthenticator
//
//  Created by Aaron Craelius on 3/31/21.
//

import SwiftUI

struct WelcomeView: View {
    var body: some View {
        VStack {
            Button(action: {}) {
                Text("Create New Account")
            }.padding(10)
            Button(action: {}) {
                Text("Sign-in to Existing Account")
            }.padding(/*@START_MENU_TOKEN@*/10/*@END_MENU_TOKEN@*/)
        }
    }
}

struct WelcomeView_Previews: PreviewProvider {
    static var previews: some View {
        WelcomeView()
    }
}
