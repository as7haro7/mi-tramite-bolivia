import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../providers/app_provider.dart';
import '../widgets/disclaimer_banner.dart';

class OnboardingScreen extends StatefulWidget {
  final VoidCallback onFinish;

  const OnboardingScreen({super.key, required this.onFinish});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  final PageController _pageController = PageController();
  int _currentPage = 0;
  String _tempCity = 'La Paz';

  final List<Map<String, dynamic>> _slides = [
    {
      'title': 'Prepara tus trámites antes de hacer fila',
      'subtitle': 'Conoce los requisitos, costos oficiales y pasos necesarios adaptados a tu caso.',
      'icon': Icons.assignment_outlined,
      'color': AppTheme.primaryBlue,
    },
    {
      'title': 'Respuestas claras con fuentes oficiales',
      'subtitle': 'Consulta a nuestro asistente conversacional con respuestas respaldadas y fecha de verificación.',
      'icon': Icons.chat_bubble_outline,
      'color': AppTheme.accentTeal,
    },
    {
      'title': 'Checklists offline y seguimiento',
      'subtitle': 'Guarda tus listas de documentos y realiza el seguimiento de tu avance sin necesidad de conexión.',
      'icon': Icons.offline_pin_outlined,
      'color': AppTheme.accentGold,
    },
  ];

  final List<String> _cities = [
    'La Paz',
    'El Alto',
    'Cochabamba',
    'Santa Cruz de la Sierra',
    'Oruro',
    'Sucre',
    'Potosí',
    'Tarija',
    'Trinidad',
    'Cobija',
    'Prefiero no indicar',
  ];

  void _finishOnboarding() {
    final provider = Provider.of<AppProvider>(context, listen: false);
    provider.setCity(_tempCity);
    provider.completeOnboarding();
    widget.onFinish();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          color: Color(0xFFF8F9FA), // Color sólido de fondo
        ),
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              children: [
                const SizedBox(height: 10),
                const DisclaimerBanner()
                    .animate()
                    .fadeIn(duration: 600.ms)
                    .slideY(begin: -0.2, end: 0),
                const SizedBox(height: 20),
                Expanded(
                  child: PageView.builder(
                    controller: _pageController,
                    onPageChanged: (index) {
                      setState(() {
                        _currentPage = index;
                      });
                    },
                    itemCount: _slides.length,
                    itemBuilder: (context, index) {
                      final slide = _slides[index];
                      return Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Container(
                            width: 120,
                            height: 120,
                            decoration: BoxDecoration(
                              color: (slide['color'] as Color).withOpacity(0.1),
                              shape: BoxShape.circle,
                            ),
                            child: Icon(
                              slide['icon'] as IconData,
                              size: 64,
                              color: slide['color'] as Color,
                            ),
                          )
                              .animate()
                              .scale(
                                duration: 600.ms,
                                curve: Curves.easeOutBack,
                              )
                              .fadeIn(duration: 400.ms),
                          const SizedBox(height: 40),
                          Text(
                            slide['title'] as String,
                            textAlign: TextAlign.center,
                            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                                  height: 1.3,
                                ),
                          )
                              .animate()
                              .fadeIn(delay: 200.ms, duration: 600.ms)
                              .slideY(begin: 0.3, end: 0, curve: Curves.easeOut),
                          const SizedBox(height: 16),
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            child: Text(
                              slide['subtitle'] as String,
                              textAlign: TextAlign.center,
                              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                                    fontSize: 16,
                                    height: 1.6,
                                    color: AppTheme.iosLightSubtext,
                                  ),
                            ),
                          )
                              .animate()
                              .fadeIn(delay: 400.ms, duration: 600.ms)
                              .slideY(begin: 0.3, end: 0, curve: Curves.easeOut),
                          if (index == _slides.length - 1) ...[
                            const SizedBox(height: 32),
                            Container(
                              decoration: BoxDecoration(
                                color: Colors.white,
                                borderRadius: BorderRadius.circular(16),
                                boxShadow: [
                                  BoxShadow(
                                    color: Colors.black.withOpacity(0.05),
                                    blurRadius: 10,
                                    offset: const Offset(0, 4),
                                  ),
                                ],
                              ),
                              child: DropdownButtonFormField<String>(
                                value: _tempCity,
                                decoration: InputDecoration(
                                  labelText: 'Tu ciudad principal (opcional)',
                                  labelStyle: TextStyle(
                                    color: AppTheme.iosLightSubtext,
                                    fontSize: 14,
                                  ),
                                  border: OutlineInputBorder(
                                    borderRadius: BorderRadius.circular(16),
                                    borderSide: BorderSide.none,
                                  ),
                                  filled: true,
                                  fillColor: Colors.white,
                                  contentPadding: const EdgeInsets.symmetric(
                                    horizontal: 20,
                                    vertical: 16,
                                  ),
                                ),
                                items: _cities.map((city) {
                                  return DropdownMenuItem(
                                    value: city,
                                    child: Text(
                                      city,
                                      style: const TextStyle(fontSize: 15),
                                    ),
                                  );
                                }).toList(),
                                onChanged: (val) {
                                  if (val != null) {
                                    setState(() {
                                      _tempCity = val;
                                    });
                                  }
                                },
                              ),
                            )
                                .animate()
                                .fadeIn(delay: 600.ms, duration: 600.ms)
                                .slideY(begin: 0.2, end: 0, curve: Curves.easeOut),
                          ],
                        ],
                      );
                    },
                  ),
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: List.generate(
                    _slides.length,
                    (index) => AnimatedContainer(
                      duration: const Duration(milliseconds: 300),
                      curve: Curves.easeInOut,
                      margin: const EdgeInsets.symmetric(horizontal: 4),
                      width: _currentPage == index ? 32 : 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: _currentPage == index
                            ? AppTheme.primaryBlue
                            : AppTheme.iosLightBorder,
                        borderRadius: BorderRadius.circular(4),
                      ),
                    ),
                  ),
                )
                    .animate()
                    .fadeIn(delay: 300.ms),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton(
                    onPressed: () {
                      if (_currentPage < _slides.length - 1) {
                        _pageController.nextPage(
                          duration: const Duration(milliseconds: 400),
                          curve: Curves.easeInOutCubic,
                        );
                      } else {
                        _finishOnboarding();
                      }
                    },
                    style: ElevatedButton.styleFrom(
                      elevation: 2,
                      shadowColor: AppTheme.primaryBlue.withOpacity(0.4),
                    ),
                    child: Text(
                      _currentPage == _slides.length - 1
                          ? 'Explorar sin registrarme'
                          : 'Siguiente',
                      style: const TextStyle(fontSize: 17),
                    ),
                  ),
                )
                    .animate()
                    .fadeIn(delay: 400.ms)
                    .slideY(begin: 0.2, end: 0),
                const SizedBox(height: 12),
                TextButton(
                  onPressed: _finishOnboarding,
                  child: const Text(
                    'Omitir e ir al inicio',
                    style: TextStyle(fontSize: 15),
                  ),
                )
                    .animate()
                    .fadeIn(delay: 500.ms),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
