# Android setup for Flutter

This guide explains how to set up the Android SDK and essential tools on Linux without installing Android Studio. We will be using the command-line tools only.

## 1. Install Java (JDK)

Android development requires the Java Development Kit (JDK). JDK 17 is generally recommended for modern Flutter and Android development.

```bash
sudo apt update
sudo apt install openjdk-17-jdk -y
```

Verify the installation:
```bash
java -version
```

## 2. Download Android Command Line Tools

You only need the Android Command Line Tools, instead of the entire Android Studio IDE package.

1. Go to the [Android Studio Downloads page](https://developer.android.com/studio#command-line-tools-only) and locate the "Command line tools only" section.
2. Download the Linux `.zip` file.
3. Alternatively, you can download it via `wget` (make sure to use the latest version link from the website):
```bash
wget https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip -O cmdline-tools.zip
```

## 3. Set Up the Android SDK Directory

The SDK tools expect a very specific directory structure (`cmdline-tools/latest/bin`).

```bash
# Create the Android SDK home directory inside your user folder
mkdir -p ~/android-sdk/cmdline-tools

# Extract the contents of the downloaded zip
unzip cmdline-tools.zip -d ~/android-sdk/cmdline-tools

# Rename the extracted 'cmdline-tools' folder to 'latest'
mv ~/android-sdk/cmdline-tools/cmdline-tools ~/android-sdk/cmdline-tools/latest
```

## 4. Set Environment Variables

You need to add the Android SDK paths to your system's environment variables. Open your shell configuration file (e.g., `~/.bashrc` or `~/.zshrc`) and add the following lines at the end:

```bash
export ANDROID_HOME=$HOME/android-sdk
export PATH=$PATH:$ANDROID_HOME/cmdline-tools/latest/bin
export PATH=$PATH:$ANDROID_HOME/platform-tools
export PATH=$PATH:$ANDROID_HOME/emulator
```

Apply the changes to your current terminal session:
```bash
source ~/.bashrc  # or source ~/.zshrc if you use zsh
```

## 5. Accept Licenses

Before installing any packages, you must accept the Android SDK licenses:

```bash
yes | sdkmanager --licenses
```

## 6. Install SDK Packages

Use the `sdkmanager` to install the essential platform-tools, build-tools, and the Android platform version you want to target (e.g., Android 34 / Android 14).

```bash
sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"
```

*(Optional)* If you plan to create and use Android Emulators locally, you also need the emulator and system images:
```bash
sdkmanager "emulator" "system-images;android-34;google_apis;x86_64"
```

## 7. Configure Flutter

Finally, point your Flutter installation to the newly created Android SDK directory:

```bash
flutter config --android-sdk ~/android-sdk
```

Run `flutter doctor` to ensure that Flutter recognizes your Android toolchain:

```bash
flutter doctor
```

If Flutter prompts you to accept any remaining Android licenses, run:
```bash
flutter doctor --android-licenses
```