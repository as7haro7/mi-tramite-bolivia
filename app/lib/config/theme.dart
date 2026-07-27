import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // Paleta de colores profesional inspirada en identidad boliviana
  static const Color primaryBlue = Color(0xFF0066CC); // Azul corporativo profundo
  static const Color primaryDarkBlue = Color(0xFF003D82); // Azul oscuro
  static const Color accentTeal = Color(0xFF00A0B0); // Teal vibrante
  static const Color accentGold = Color(0xFFFDB913); // Dorado boliviano
  
  // iOS HIG System Colors (mantenemos para compatibilidad)
  static const Color tintColor = primaryBlue;
  static const Color systemTeal = accentTeal;
  static const Color systemGreen = Color(0xFF34C759);
  static const Color systemOrange = Color(0xFFFF9500);
  static const Color systemRed = Color(0xFFFF3B30);

  // Aliases for compatibility
  static const Color primaryViolet = primaryBlue;
  static const Color iosBlue = primaryBlue;
  static const Color iosTeal = systemTeal;
  static const Color secondaryTeal = systemTeal;
  static const Color successGreen = systemGreen;
  static const Color alertRed = systemRed;
  static const Color accentAmber = accentGold;
  static const Color warningOrange = systemOrange;

  // Backgrounds y superficies - Tema claro profesional
  static const Color iosLightBg = Color(0xFFF8F9FA); // Gris muy claro
  static const Color iosLightCard = Color(0xFFFFFFFF);
  static const Color iosLightText = Color(0xFF1A1A1A); // Negro suave
  static const Color iosLightSubtext = Color(0xFF6B7280); // Gris medio
  static const Color iosLightBorder = Color(0xFFE5E7EB);
  static const Color iosLightSearchBg = Color(0xFFF3F4F6);
  
  // Superficies adicionales para elevación
  static const Color surfaceElevated = Color(0xFFFFFFFF);
  static const Color surfaceOverlay = Color(0xFFF9FAFB);

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

  // iOS Subtle Gradients - Actualizados con colores profesionales
  // Nota: Se prefiere usar colores sólidos con efectos blur para diseño moderno
  static const LinearGradient primaryGradient = LinearGradient(
    colors: [Color(0xFF0066CC), Color(0xFF0066CC)], // Color sólido
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  static const LinearGradient accentGradient = LinearGradient(
    colors: [Color(0xFF00A0B0), Color(0xFF00A0B0)], // Color sólido
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  static const LinearGradient premiumGradient = LinearGradient(
    colors: [Color(0xFFFDB913), Color(0xFFFDB913)], // Color sólido
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
  
  static const LinearGradient subtleGradient = LinearGradient(
    colors: [Color(0xFFF8F9FA), Color(0xFFF8F9FA)], // Color sólido
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  // Light Theme (Profesional y Moderno)
  static ThemeData get lightTheme {
    final textTheme = GoogleFonts.interTextTheme(ThemeData.light().textTheme);

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: ColorScheme.fromSeed(
        seedColor: primaryBlue,
        primary: primaryBlue,
        secondary: accentTeal,
        tertiary: accentGold,
        surface: iosLightBg,
        onSurface: iosLightText,
        error: systemRed,
        brightness: Brightness.light,
      ),
      scaffoldBackgroundColor: iosLightBg,
      textTheme: textTheme.copyWith(
        headlineLarge: textTheme.headlineLarge?.copyWith(
          fontSize: 32,
          fontWeight: FontWeight.w800,
          color: iosLightText,
          letterSpacing: -0.8,
          height: 1.2,
        ),
        headlineMedium: textTheme.headlineMedium?.copyWith(
          fontSize: 26,
          fontWeight: FontWeight.w700,
          color: iosLightText,
          letterSpacing: -0.6,
          height: 1.3,
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
          letterSpacing: -0.2,
        ),
        bodyLarge: textTheme.bodyLarge?.copyWith(
          fontSize: 16,
          color: iosLightText,
          height: 1.5,
        ),
        bodyMedium: textTheme.bodyMedium?.copyWith(
          fontSize: 14,
          color: iosLightSubtext,
          height: 1.5,
        ),
        labelLarge: textTheme.labelLarge?.copyWith(
          fontSize: 15,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.1,
        ),
      ),
      cardTheme: CardThemeData(
        color: iosLightCard,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: BorderSide(color: iosLightBorder.withOpacity(0.3), width: 1),
        ),
      ),
      appBarTheme: AppBarTheme(
        backgroundColor: iosLightCard,
        elevation: 0,
        centerTitle: true,
        scrolledUnderElevation: 0,
        surfaceTintColor: Colors.transparent,
        iconTheme: const IconThemeData(color: primaryBlue, size: 22),
        titleTextStyle: TextStyle(
          color: iosLightText,
          fontSize: 18,
          fontWeight: FontWeight.w700,
          letterSpacing: -0.3,
          fontFamily: GoogleFonts.inter().fontFamily,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primaryBlue,
          foregroundColor: Colors.white,
          elevation: 0,
          shadowColor: primaryBlue.withOpacity(0.3),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          textStyle: const TextStyle(
            fontWeight: FontWeight.w700,
            fontSize: 16,
            letterSpacing: 0.2,
          ),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: primaryBlue,
          textStyle: const TextStyle(
            fontWeight: FontWeight.w600,
            fontSize: 15,
            letterSpacing: 0.1,
          ),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: iosLightSearchBg,
        labelStyle: TextStyle(
          color: iosLightText,
          fontSize: 13,
          fontWeight: FontWeight.w500,
        ),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: iosLightSearchBg,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: primaryBlue, width: 2),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        hintStyle: TextStyle(color: iosLightSubtext, fontSize: 15),
      ),
      dividerTheme: DividerThemeData(
        color: iosLightBorder.withOpacity(0.5),
        thickness: 1,
        space: 1,
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
