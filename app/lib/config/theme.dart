import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // iOS HIG System Tint Color
  static const Color tintColor = Color(0xFF007AFF); // System Blue
  static const Color systemTeal = Color(0xFF30B0C7);
  static const Color systemGreen = Color(0xFF34C759);
  static const Color systemOrange = Color(0xFFFF9500);
  static const Color systemRed = Color(0xFFFF3B30);

  // Aliases for compatibility
  static const Color primaryBlue = tintColor;
  static const Color primaryViolet = tintColor;
  static const Color iosBlue = tintColor;
  static const Color iosTeal = systemTeal;
  static const Color secondaryTeal = systemTeal;
  static const Color successGreen = systemGreen;
  static const Color alertRed = systemRed;
  static const Color accentAmber = systemOrange;

  // iOS Light System Colors
  static const Color iosLightBg = Color(0xFFF2F2F7); // Grouped Table Background
  static const Color iosLightCard = Color(0xFFFFFFFF);
  static const Color iosLightText = Color(0xFF000000);
  static const Color iosLightSubtext = Color(0xFF8E8E93);
  static const Color iosLightBorder = Color(0xFFC6C6C8);
  static const Color iosLightSearchBg = Color(0xFFE3E3E8);

  // iOS Dark System Colors
  static const Color iosDarkBg = Color(0xFF000000); // System Pure Black
  static const Color iosDarkCard = Color(0xFF1C1C1E); // Grouped Elevated Dark
  static const Color iosDarkText = Color(0xFFFFFFFF);
  static const Color iosDarkSubtext = Color(0xFF8E8E93);
  static const Color iosDarkBorder = Color(0xFF38383A);
  static const Color iosDarkSearchBg = Color(0xFF1C1C1E);

  // Aliases
  static const Color lightBorder = iosLightBorder;
  static const Color darkBorder = iosDarkBorder;
  static const Color textHeading = iosLightText;
  static const Color textBody = iosLightSubtext;
  static const Color textMuted = Color(0xFF8E8E93);

  // iOS Subtle Gradients
  static const LinearGradient primaryGradient = LinearGradient(
    colors: [Color(0xFF007AFF), Color(0xFF0051A8)],
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  static const LinearGradient premiumGradient = LinearGradient(
    colors: [Color(0xFF5856D6), Color(0xFFAF52DE)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  // Light Theme (Strict Apple HIG)
  static ThemeData get lightTheme {
    final textTheme = GoogleFonts.interTextTheme(ThemeData.light().textTheme);

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: ColorScheme.fromSeed(
        seedColor: tintColor,
        primary: tintColor,
        secondary: systemTeal,
        surface: iosLightBg,
        onSurface: iosLightText,
        error: systemRed,
      ),
      scaffoldBackgroundColor: iosLightBg,
      textTheme: textTheme.copyWith(
        headlineLarge: textTheme.headlineLarge?.copyWith(
          fontSize: 34,
          fontWeight: FontWeight.bold,
          color: iosLightText,
          letterSpacing: -0.4,
        ),
        headlineMedium: textTheme.headlineMedium?.copyWith(
          fontSize: 28,
          fontWeight: FontWeight.bold,
          color: iosLightText,
          letterSpacing: -0.4,
        ),
        titleLarge: textTheme.titleLarge?.copyWith(
          fontSize: 20,
          fontWeight: FontWeight.w700,
          color: iosLightText,
          letterSpacing: -0.4,
        ),
        titleMedium: textTheme.titleMedium?.copyWith(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          color: iosLightText,
        ),
        bodyLarge: textTheme.bodyLarge?.copyWith(
          fontSize: 16,
          color: iosLightText,
        ),
        bodyMedium: textTheme.bodyMedium?.copyWith(
          fontSize: 14,
          color: iosLightSubtext,
        ),
      ),
      cardTheme: CardThemeData(
        color: iosLightCard,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
        scrolledUnderElevation: 0,
        iconTheme: IconThemeData(color: tintColor),
        titleTextStyle: TextStyle(
          color: iosLightText,
          fontSize: 17,
          fontWeight: FontWeight.w600,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: tintColor,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
          textStyle: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: tintColor,
          textStyle: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
        ),
      ),
    );
  }

  // Dark Theme (Strict Apple HIG Dark)
  static ThemeData get darkTheme {
    final textTheme = GoogleFonts.interTextTheme(ThemeData.dark().textTheme);

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: ColorScheme.fromSeed(
        seedColor: tintColor,
        brightness: Brightness.dark,
        primary: tintColor,
        secondary: systemTeal,
        surface: iosDarkBg,
        onSurface: iosDarkText,
        error: systemRed,
      ),
      scaffoldBackgroundColor: iosDarkBg,
      textTheme: textTheme.copyWith(
        headlineLarge: textTheme.headlineLarge?.copyWith(
          fontSize: 34,
          fontWeight: FontWeight.bold,
          color: iosDarkText,
          letterSpacing: -0.4,
        ),
        headlineMedium: textTheme.headlineMedium?.copyWith(
          fontSize: 28,
          fontWeight: FontWeight.bold,
          color: iosDarkText,
          letterSpacing: -0.4,
        ),
        titleLarge: textTheme.titleLarge?.copyWith(
          fontSize: 20,
          fontWeight: FontWeight.w700,
          color: iosDarkText,
          letterSpacing: -0.4,
        ),
        titleMedium: textTheme.titleMedium?.copyWith(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          color: iosDarkText,
        ),
        bodyLarge: textTheme.bodyLarge?.copyWith(
          fontSize: 16,
          color: iosDarkText,
        ),
        bodyMedium: textTheme.bodyMedium?.copyWith(
          fontSize: 14,
          color: iosDarkSubtext,
        ),
      ),
      cardTheme: CardThemeData(
        color: iosDarkCard,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
        scrolledUnderElevation: 0,
        iconTheme: IconThemeData(color: tintColor),
        titleTextStyle: TextStyle(
          color: iosDarkText,
          fontSize: 17,
          fontWeight: FontWeight.w600,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: tintColor,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
          textStyle: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16),
        ),
      ),
    );
  }
}
