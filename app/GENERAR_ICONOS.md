# 🎨 Generar Iconos de la Aplicación

## Pasos para Generar los Iconos

### 1. Verificar que el logo esté en su lugar
El logo ya está configurado en `assets/images/logo.png`

### 2. Ejecutar el generador de iconos

```bash
# Navega al directorio de la app
cd app

# Instala las dependencias (si aún no lo hiciste)
flutter pub get

# Genera los iconos automáticamente
dart run flutter_launcher_icons
```

### 3. Verificar la generación

El comando generará iconos para:
- ✅ **Android**: Todos los tamaños (mipmap)
- ✅ **iOS**: Todos los tamaños (Assets.xcassets)
- ✅ **Web**: favicon y apple-touch-icon

### 4. Resultado Esperado

Verás un output similar a:
```
Creating icons...
✓ Android
  ✓ Creating mipmap-mdpi
  ✓ Creating mipmap-hdpi
  ✓ Creating mipmap-xhdpi
  ✓ Creating mipmap-xxhdpi
  ✓ Creating mipmap-xxxhdpi
✓ iOS
  ✓ Creating AppIcon.appiconset
✓ Web
  ✓ Creating favicon.png
  ✓ Creating apple-touch-icon.png
```

## Configuración Actual

El archivo `pubspec.yaml` está configurado con:

```yaml
flutter_launcher_icons:
  android: true
  ios: true
  web:
    generate: true
    image_path: "assets/images/logo.png"
    background_color: "#0066CC"
  image_path: "assets/images/logo.png"
  adaptive_icon_background: "#0066CC"
  adaptive_icon_foreground: "assets/images/logo.png"
```

### Características:
- **Fondo azul corporativo**: `#0066CC`
- **Logo en primer plano**: `logo.png`
- **Iconos adaptativos** para Android 8.0+
- **Favicon** para web

## Verificar los Iconos

### Android
Los iconos se generan en:
```
android/app/src/main/res/mipmap-*/ic_launcher.png
```

### iOS
Los iconos se generan en:
```
ios/Runner/Assets.xcassets/AppIcon.appiconset/
```

### Web
Los iconos se generan en:
```
web/icons/
web/favicon.png
```

## Personalización Adicional

Si deseas cambiar el color de fondo o usar otro logo:

1. Edita `pubspec.yaml`:
```yaml
flutter_launcher_icons:
  adaptive_icon_background: "#TU_COLOR_AQUI"
  image_path: "ruta/a/tu/logo.png"
```

2. Vuelve a ejecutar:
```bash
dart run flutter_launcher_icons
```

## Troubleshooting

### Error: "No pubspec.yaml found"
Asegúrate de estar en el directorio `app/`:
```bash
cd app
```

### Error: "Image not found"
Verifica que el logo exista en:
```bash
ls assets/images/logo.png
```

### Los iconos no se actualizan
1. Limpia el proyecto:
```bash
flutter clean
flutter pub get
```

2. Vuelve a generar:
```bash
dart run flutter_launcher_icons
```

3. Reconstruye la app:
```bash
flutter run
```

## Probar los Iconos

### Android
```bash
flutter run -d android
```
El icono aparecerá en el launcher/home screen

### iOS
```bash
flutter run -d ios
```
El icono aparecerá en el home screen

### Web
```bash
flutter run -d chrome
```
El favicon aparecerá en la pestaña del navegador

## Resultado Final

Una vez generados, tu app tendrá:
- ✅ Icono corporativo con el logo
- ✅ Fondo azul corporativo `#0066CC`
- ✅ Iconos adaptativos para Android moderno
- ✅ Compatibilidad con todas las plataformas

---

**Nota**: Los iconos se generan automáticamente en todos los tamaños requeridos por cada plataforma.
