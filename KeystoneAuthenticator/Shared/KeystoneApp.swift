//
//  KeystoneApp.swift
//  Shared
//
//  Created by Aaron Craelius on 3/30/21.
//

import SwiftUI

@main
struct KeystoneApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
        }.onChange(of: scenePhase) { _ in
            persistenceController.save()
        }
    }
    
    let persistenceController = PersistenceController.shared
    
    @Environment(\.scenePhase) var scenePhase
}
