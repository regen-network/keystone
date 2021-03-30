//
//  SettingsView.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import SwiftUI

struct SettingsView: View {
    var body: some View {
        NavigationView {
            List {
                Label("Security", systemImage: "lock")
                Label("Fees", systemImage: "banknote")
                Label("Notifications", systemImage: "exclamationmark.bubble")
                Label("Logut", systemImage: "arrow.down.left.circle")
            }
        }
    }
}

struct SettingsView_Previews: PreviewProvider {
    static var previews: some View {
        SettingsView()
    }
}
