import 'package:flutter/material.dart';
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

  final List<Map<String, String>> _slides = [
    {
      'title': 'Prepara tus trámites antes de hacer fila',
      'subtitle': 'Conoce los requisitos, costos oficiales y pasos necesarios adaptados a tu caso.',
      'icon': '📋',
    },
    {
      'title': 'Respuestas claras con fuentes oficiales',
      'subtitle': 'Consulta a nuestro asistente conversacional con respuestas respaldadas y fecha de verificación.',
      'icon': '💬',
    },
    {
      'title': 'Checklists offline y seguimiento',
      'subtitle': 'Guarda tus listas de documentos y realiza el seguimiento de tu avance sin necesidad de conexión.',
      'icon': '📱',
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
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            children: [
              const SizedBox(height: 10),
              const DisclaimerBanner(),
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
                        Text(
                          slide['icon']!,
                          style: const TextStyle(fontSize: 70),
                        ),
                        const SizedBox(height: 24),
                        Text(
                          slide['title']!,
                          textAlign: TextAlign.center,
                          style: Theme.of(context).textTheme.headlineMedium,
                        ),
                        const SizedBox(height: 14),
                        Text(
                          slide['subtitle']!,
                          textAlign: TextAlign.center,
                          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                fontSize: 15,
                                height: 1.4,
                              ),
                        ),
                        if (index == _slides.length - 1) ...[
                          const SizedBox(height: 24),
                          DropdownButtonFormField<String>(
                            initialValue: _tempCity,
                            decoration: InputDecoration(
                              labelText: 'Tu ciudad principal (opcional)',
                              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                            ),
                            items: _cities.map((city) {
                              return DropdownMenuItem(value: city, child: Text(city));
                            }).toList(),
                            onChanged: (val) {
                              if (val != null) {
                                setState(() {
                                  _tempCity = val;
                                });
                              }
                            },
                          ),
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
                  (index) => Container(
                    margin: const EdgeInsets.symmetric(horizontal: 4),
                    width: _currentPage == index ? 24 : 8,
                    height: 8,
                    decoration: BoxDecoration(
                      color: _currentPage == index ? AppTheme.primaryBlue : Colors.grey.shade400,
                      borderRadius: BorderRadius.circular(4),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    if (_currentPage < _slides.length - 1) {
                      _pageController.nextPage(
                        duration: const Duration(milliseconds: 300),
                        curve: Curves.easeInOut,
                      );
                    } else {
                      _finishOnboarding();
                    }
                  },
                  child: Text(_currentPage == _slides.length - 1 ? 'Explorar sin registrarme' : 'Siguiente'),
                ),
              ),
              const SizedBox(height: 10),
              TextButton(
                onPressed: _finishOnboarding,
                child: const Text('Omitir e ir al inicio'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
