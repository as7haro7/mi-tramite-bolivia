import 'dart:ui';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../config/theme.dart';

class BottomNavBar extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTap;

  const BottomNavBar({
    super.key,
    required this.currentIndex,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
        child: Container(
          decoration: const BoxDecoration(
            color: Color(0xCCF9F9F9),
            border: Border(
              top: BorderSide(
                color: Color(0xFFC6C6C8),
                width: 0.5,
              ),
            ),
          ),
          child: CupertinoTabBar(
            currentIndex: currentIndex,
            onTap: (index) {
              HapticFeedback.selectionClick();
              onTap(index);
            },
            backgroundColor: Colors.transparent,
            activeColor: AppTheme.tintColor,
            inactiveColor: const Color(0xFF8E8E93),
            iconSize: 24,
            items: const [
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.house_fill),
                label: 'Inicio',
              ),
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.bookmark_fill),
                label: 'Guardados',
              ),
              BottomNavigationBarItem(
                icon: Icon(CupertinoIcons.person_crop_circle_fill),
                label: 'Perfil',
              ),
            ],
          ),
        ),
      ),
    );
  }
}
