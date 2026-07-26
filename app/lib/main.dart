import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'config/theme.dart';
import 'models/tramite.dart';
import 'providers/app_provider.dart';
import 'providers/assistant_provider.dart';
import 'providers/checklist_provider.dart';
import 'providers/favorites_provider.dart';
import 'screens/assistant_screen.dart';
import 'screens/home_screen.dart';
import 'screens/onboarding_screen.dart';
import 'screens/premium_plan_screen.dart';
import 'screens/profile_screen.dart';
import 'screens/saved_screen.dart';
import 'screens/search_results_screen.dart';
import 'screens/tramite_detail_screen.dart';
import 'widgets/bottom_nav_bar.dart';
import 'widgets/floating_assistant_button.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.dark,
    ),
  );
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AppProvider()),
        ChangeNotifierProvider(create: (_) => FavoritesProvider()),
        ChangeNotifierProvider(create: (_) => ChecklistProvider()),
        ChangeNotifierProvider(create: (_) => AssistantProvider()),
      ],
      child: const MiTramiteApp(),
    ),
  );
}

class MiTramiteApp extends StatelessWidget {
  const MiTramiteApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Mi Trámite Bolivia',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.lightTheme,
      themeMode: ThemeMode.light,
      home: const MainNavigationWrapper(),
    );
  }
}

class MainNavigationWrapper extends StatefulWidget {
  const MainNavigationWrapper({super.key});

  @override
  State<MainNavigationWrapper> createState() => _MainNavigationWrapperState();
}

class _MainNavigationWrapperState extends State<MainNavigationWrapper> {
  int _currentIndex = 0;
  bool _showingOnboarding = false;

  @override
  void initState() {
    super.initState();
    _checkOnboarding();
  }

  void _checkOnboarding() {
    final appProvider = Provider.of<AppProvider>(context, listen: false);
    if (!appProvider.onboardingCompleted) {
      setState(() {
        _showingOnboarding = true;
      });
    }
  }

  void _navigateToTramiteDetail(Tramite tramite) {
    HapticFeedback.lightImpact();
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (_) => TramiteDetailScreen(
          tramite: tramite,
          onChecklistCreated: () {
            setState(() {
              _currentIndex = 1; // Tab 1 = Mis Guardados / Checklists
            });
          },
        ),
      ),
    );
  }

  void _navigateToSearchResults(String query) {
    HapticFeedback.lightImpact();
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (_) => SearchResultsScreen(
          initialQuery: query,
          onSelectTramite: _navigateToTramiteDetail,
          onOpenAssistant: () {
            Navigator.pop(context);
            _openAssistantModal();
          },
        ),
      ),
    );
  }

  void _openAssistantModal() {
    HapticFeedback.lightImpact();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => FractionallySizedBox(
        heightFactor: 0.88,
        child: Container(
          decoration: const BoxDecoration(
            color: AppTheme.iosLightBg,
            borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
            boxShadow: [
              BoxShadow(color: Colors.black12, blurRadius: 20, offset: Offset(0, -4)),
            ],
          ),
          child: Column(
            children: [
              const SizedBox(height: 10),
              Container(
                width: 36,
                height: 5,
                decoration: BoxDecoration(
                  color: Colors.grey.shade400,
                  borderRadius: BorderRadius.circular(2.5),
                ),
              ),
              const SizedBox(height: 6),
              const Expanded(child: AssistantScreen(isBottomSheet: true)),
            ],
          ),
        ),
      ),
    );
  }

  void _navigateToPremiumPlan() {
    HapticFeedback.lightImpact();
    Navigator.push(
      context,
      MaterialPageRoute(builder: (_) => const PremiumPlanScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    final appProvider = Provider.of<AppProvider>(context);

    if (_showingOnboarding && !appProvider.onboardingCompleted) {
      return OnboardingScreen(
        onFinish: () {
          setState(() {
            _showingOnboarding = false;
          });
        },
      );
    }

    final pages = [
      HomeScreen(
        onSearchSubmitted: _navigateToSearchResults,
        onOpenAssistant: _openAssistantModal,
        onSelectTramite: _navigateToTramiteDetail,
        onOpenPremium: _navigateToPremiumPlan,
      ),
      SavedScreen(
        onSelectTramite: _navigateToTramiteDetail,
        onOpenPremium: _navigateToPremiumPlan,
      ),
      ProfileScreen(
        onOpenPremium: _navigateToPremiumPlan,
      ),
    ];

    return Scaffold(
      extendBody: true,
      body: IndexedStack(
        index: _currentIndex,
        children: pages,
      ),
      floatingActionButton: const FloatingAssistantButton(),
      floatingActionButtonLocation: FloatingActionButtonLocation.endFloat,
      bottomNavigationBar: BottomNavBar(
        currentIndex: _currentIndex,
        onTap: (index) => setState(() => _currentIndex = index),
      ),
    );
  }
}
