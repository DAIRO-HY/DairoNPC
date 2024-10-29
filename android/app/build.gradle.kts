plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "cn.dairo.npc"
    compileSdk = 34

    defaultConfig {
        applicationId = "cn.dairo.npc"
        minSdk = 21
        targetSdk = 34
        versionCode = 1
        versionName = "1.0"

        vectorDrawables {
            useSupportLibrary = true
        }
    }
    signingConfigs {
        create("releaseConfig") {
            storeFile = File(rootProject.projectDir, "dairo-npc.jks")

            //Ths********1
            storePassword = System.getenv("DAIRO-NPC-JKS-PASSWORD")
            keyAlias = "dairo-npc"
            keyPassword = System.getenv("DAIRO-NPC-JKS-PASSWORD")
        }
    }

    buildTypes {
        release {
            signingConfig = signingConfigs.getByName("releaseConfig")
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_1_8
        targetCompatibility = JavaVersion.VERSION_1_8
    }
    kotlinOptions {
        jvmTarget = "1.8"
    }
    packaging {
        resources {
            excludes += "/META-INF/{AL2.0,LGPL2.1}"
        }
    }
}

dependencies {
    implementation("androidx.core:core-ktx:1.13.1")
    implementation(files("libs/dairo-npc.aar"))
}