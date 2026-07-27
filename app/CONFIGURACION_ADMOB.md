# Configuración de Google AdMob

Esta guía te ayudará a configurar correctamente Google AdMob en la aplicación Flutter.

## 1. Crear Cuenta en Google AdMob

1. Ve a [https://admob.google.com](https://admob.google.com)
2. Inicia sesión con tu cuenta de Google
3. Sigue el proceso de registro
4. Acepta los términos y condiciones

## 2. Crear una Aplicación en AdMob

### Crear App
1. En el panel de AdMob, haz clic en "Aplicaciones" > "Agregar aplicación"
2. Selecciona la plataforma (Android/iOS)
3. Si aún no está publicada, selecciona "No"
4. Ingresa el nombre: "Mi Trámite Bolivia"
5. Haz clic en "Agregar"

### Obtener el Application ID
- **Android**: Algo como `ca-app-pub-3940256099942544~3347511713`
- **iOS**: Algo como `ca-app-pub-3940256099942544~1458002511`

## 3. Crear Unidades de Anuncios

### Para Banner (necesario para la app)
1. En tu aplicación, ve a "Unidades de anuncios"
2. Haz clic en "Comenzar" o "Crear unidad de anuncios"
3. Selecciona "Banner"
4. Configura:
   - Nombre: "Banner Principal"
   - Formato: Estándar (320x50)
5. Haz clic en "Crear unidad de anuncios"
6. **Guarda el ID** (ejemplo: `ca-app-pub-3940256099942544/6300978111`)

Repite este proceso para Android e iOS si son diferentes.

## 4. Configurar Android

### Paso 1: Editar AndroidManifest.xml
Ubicación: `android/app/src/main/AndroidManifest.xml`

Agrega dentro de la etiqueta `<application>` (después de `android:label`):

```xml
<application
    android:label="Mi Trámite Bolivia"
    android:name="${applicationName}"
    android:icon="@mipmap/ic_launcher">
    
    <!-- AGREGAR ESTA LÍNEA -->
    <meta-data
        android:name="com.google.android.gms.ads.APPLICATION_ID"
        android:value="ca-app-pub-XXXXXXXXXXXXXXXX~YYYYYYYYYY"/>
    
    <!-- Resto del contenido -->
    <activity>
        ...
    </activity>
</application>
```

**Reemplaza** `ca-app-pub-XXXXXXXXXXXXXXXX~YYYYYYYYYY` con tu Application ID real de Android.

### Paso 2: Verificar permisos de Internet
Asegúrate de que existe en `AndroidManifest.xml` (fuera de `<application>`):

```xml
<uses-permission android:name="android.permission.INTERNET"/>
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE"/>
```

## 5. Configurar iOS

### Paso 1: Editar Info.plist
Ubicación: `ios/Runner/Info.plist`

Agrega antes de la etiqueta final `</dict>`:

```xml
<dict>
    <!-- Configuración existente -->
    
    <!-- AGREGAR ESTAS LÍNEAS -->
    <key>GADApplicationIdentifier</key>
    <string>ca-app-pub-XXXXXXXXXXXXXXXX~YYYYYYYYYY</string>
    
    <key>SKAdNetworkItems</key>
    <array>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>cstr6suwn9.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>4fzdc2evr5.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>4pfyvq9l8r.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>2fnua5tdw4.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>ydx93a7ass.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>5a6flpkh64.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>p78axxw29g.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>v72qych5uu.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>ludvb6z3bs.skadnetwork</string>
        </dict>
        <dict>
            <key>SKAdNetworkIdentifier</key>
            <string>cp8zw746q7.skadnetwork</string>
        </dict>
    </array>
    
    <!-- App Transport Security -->
    <key>NSAppTransportSecurity</key>
    <dict>
        <key>NSAllowsArbitraryLoads</key>
        <true/>
    </dict>
</dict>
```

**Reemplaza** `ca-app-pub-XXXXXXXXXXXXXXXX~YYYYYYYYYY` con tu Application ID real de iOS.

### Paso 2: Actualizar Podfile (si es necesario)
Ubicación: `ios/Podfile`

Asegúrate de tener:
```ruby
platform :ios, '12.0'  # O superior
```

## 6. Actualizar ad_config.dart

Ubicación: `lib/config/ad_config.dart`

```dart
class AdConfig {
  // IDs de PRODUCCIÓN
  static const String _prodBannerAdUnitIdAndroid = 'ca-app-pub-XXXXXXXX/YYYYYY';
  static const String _prodBannerAdUnitIdIOS = 'ca-app-pub-XXXXXXXX/YYYYYY';

  static String get bannerAdUnitId {
    if (Platform.isAndroid) {
      return _prodBannerAdUnitIdAndroid; // Usar en producción
    } else if (Platform.isIOS) {
      return _prodBannerAdUnitIdIOS; // Usar en producción
    }
    throw UnsupportedError('Plataforma no soportada');
  }
}
```

## 7. Testing y Validación

### Usar IDs de Prueba (Desarrollo)
Durante el desarrollo, **SIEMPRE** usa los IDs de prueba de Google:

**Android Banner Test ID**: `ca-app-pub-3940256099942544/6300978111`
**iOS Banner Test ID**: `ca-app-pub-3940256099942544/2934735716`

### Probar la Integración
```bash
# Limpiar y reconstruir
flutter clean
flutter pub get

# Ejecutar en Android
flutter run

# Ejecutar en iOS
flutter run -d ios
```

### Verificar que Funciona
1. Los anuncios deben aparecer en:
   - Pantalla de inicio (después del CTA Premium)
   - Detalle de trámite (pestaña Requisitos)
   - Pantalla de guardados (ambas pestañas)

2. Deberías ver anuncios de prueba con texto "Test Ad"

## 8. Publicación en Producción

### Antes de Publicar
- [ ] Reemplazar IDs de prueba con IDs reales
- [ ] Verificar que `shouldShowAds` está configurado correctamente
- [ ] Probar en dispositivos reales
- [ ] Revisar políticas de AdMob
- [ ] Configurar medición en AdMob

### Checklist de Producción
```dart
// En ad_config.dart
static String get bannerAdUnitId {
  if (Platform.isAndroid) {
    return _prodBannerAdUnitIdAndroid; // ✅ Debe ser ID real
  } else if (Platform.isIOS) {
    return _prodBannerAdUnitIdIOS; // ✅ Debe ser ID real
  }
}
```

## 9. Monitoreo y Optimización

### Panel de AdMob
1. Monitorea impresiones diarias
2. Revisa tasa de clics (CTR)
3. Verifica ingresos estimados
4. Optimiza ubicaciones según rendimiento

### Métricas Importantes
- **Impresiones**: Cantidad de veces que se muestra un anuncio
- **CTR (Click-Through Rate)**: Porcentaje de clics sobre impresiones
- **eCPM**: Ingresos estimados por 1000 impresiones
- **Fill Rate**: Porcentaje de solicitudes con anuncio mostrado

## 10. Troubleshooting

### Los anuncios no se muestran
1. Verifica que Application ID esté configurado correctamente
2. Revisa los logs: `flutter run --verbose`
3. Asegúrate de tener conexión a internet
4. Espera unos minutos (los anuncios pueden tardar en cargar)
5. Verifica que los IDs de prueba sean correctos

### Error "Invalid Ad Unit ID"
- El ID de unidad de anuncio es incorrecto
- Verifica que copiaste el ID completo desde AdMob
- Asegúrate de usar el ID correcto para cada plataforma

### Anuncios en producción pero no generan ingresos
- Estás usando IDs de prueba (reemplázalos con IDs reales)
- La aplicación aún no está aprobada en AdMob (puede tomar horas/días)
- Necesitas más tráfico para generar ingresos significativos

## Recursos Adicionales

- [Google AdMob](https://admob.google.com)
- [Documentación Flutter google_mobile_ads](https://pub.dev/packages/google_mobile_ads)
- [Políticas de AdMob](https://support.google.com/admob/answer/6128543)
- [Mejores prácticas de implementación](https://support.google.com/admob/answer/6128877)

---

**Nota**: Esta es una versión DEMO. Los anuncios actuales usan IDs de prueba y no generan ingresos reales.
