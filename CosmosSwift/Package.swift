// swift-tools-version:5.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "CosmosSwift",
    products: [
        // Products define the executables and libraries a package produces, and make them visible to other packages.
        .library(
            name: "CosmosSwift",
            targets: ["CosmosSwift"]),
    ],
    dependencies: [
        .package(name: "SwiftProtobuf", url: "https://github.com/apple/swift-protobuf.git", from: "1.6.0"),
        .package(url: "https://github.com/grpc/grpc-swift.git", from: "1.0.0-alpha.21")
    ],
    targets: [
        // Targets are the basic building blocks of a package. A target can define a module or a test suite.
        // Targets can depend on other targets in this package, and on products in packages this package depends on.
        .target(
            name: "CosmosSwift",
            dependencies: ["SwiftProtobuf", .product(name: "GRPC", package: "grpc-swift")]),
        .testTarget(
            name: "CosmosSwiftTests",
            dependencies: ["CosmosSwift"]),
    ]
)
