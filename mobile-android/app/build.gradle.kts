plugins { id("com.android.application"); id("org.jetbrains.kotlin.android"); id("com.google.devtools.ksp") }

android { namespace = "mx.edu.accesoindustrial"; compileSdk = 35
    defaultConfig { applicationId = "mx.edu.accesoindustrial"; minSdk = 26; targetSdk = 35; versionCode = 1; versionName = "1.0" }
    buildFeatures { viewBinding = true; buildConfig = true }
    // AGP and KSP must target the same bytecode version. Gradle can run on the
    // installed JDK 21 while the Android app is emitted as Java 17 bytecode.
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
    kotlinOptions { jvmTarget = "17" }
}
dependencies {
    implementation("androidx.core:core-ktx:1.15.0"); implementation("androidx.appcompat:appcompat:1.7.0"); implementation("com.google.android.material:material:1.12.0")
    implementation("androidx.lifecycle:lifecycle-viewmodel-ktx:2.8.7"); implementation("androidx.lifecycle:lifecycle-livedata-ktx:2.8.7")
    implementation("androidx.camera:camera-core:1.4.1"); implementation("androidx.camera:camera-camera2:1.4.1"); implementation("androidx.camera:camera-lifecycle:1.4.1"); implementation("androidx.camera:camera-view:1.4.1")
    implementation("com.journeyapps:zxing-android-embedded:4.3.0")
    implementation("com.squareup.retrofit2:retrofit:2.11.0"); implementation("com.squareup.retrofit2:converter-gson:2.11.0"); implementation("com.squareup.okhttp3:logging-interceptor:4.12.0")
    implementation("androidx.room:room-runtime:2.6.1"); implementation("androidx.room:room-ktx:2.6.1"); ksp("androidx.room:room-compiler:2.6.1")
    implementation("androidx.work:work-runtime-ktx:2.10.0")
}
