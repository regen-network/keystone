//
//  ApproveTxView.swift
//  Keystone
//
//  Created by Aaron Craelius on 3/30/21.
//

import SwiftUI

struct ApproveTxView: View {
    let tx: TxInfo
    
    var body: some View {
        VStack {
            Text(tx.summary)
                .fontWeight(.bold)
            Button(action: {
                
            }) {
                Text("Approve")
            }.buttonStyle(ApproveButtonStyle(bgColor: Color.green))
            Button(action: {
                
            }) {
                Text("Reject")
            }.buttonStyle(ApproveButtonStyle(bgColor: Color.red))
        }
    }
}

struct ApproveButtonStyle: ButtonStyle {
    let bgColor: Color
    public func makeBody(configuration: Self.Configuration) -> some View {
        configuration.label
            .font(Font.body.weight(.medium))
            .padding(.vertical, 12)
            .foregroundColor(Color.white)
            .frame(maxWidth: .infinity)
            .background(
                    RoundedRectangle(cornerRadius: 14.0, style: .continuous)
                    .fill(bgColor)
            )
            .opacity(configuration.isPressed ? 0.4 : 1.0)
    }
}

struct ApproveTxView_Previews: PreviewProvider {
    static var previews: some View {
        ApproveTxView(tx: TxInfo(summary: "Send 100REGEN from regen1abcdecxvad49sdf2 to regen1sdgkhwfeiaflknewio"))
    }
}
