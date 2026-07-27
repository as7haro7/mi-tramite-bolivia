import 'package:flutter/material.dart';
import '../config/theme.dart';

/// Widget que muestra una animación de "typing" (escribiendo)
/// para indicar que la IA está procesando una respuesta.
class TypingIndicator extends StatefulWidget {
  const TypingIndicator({super.key});

  @override
  State<TypingIndicator> createState() => _TypingIndicatorState();
}

class _TypingIndicatorState extends State<TypingIndicator>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1400),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Avatar del asistente
          Container(
            width: 32,
            height: 32,
            decoration: const BoxDecoration(
              gradient: AppTheme.primaryGradient,
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.auto_awesome_rounded,
              color: Colors.white,
              size: 16,
            ),
          ),
          const SizedBox(width: 8),

          // Bubble con animación de typing
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: BoxDecoration(
              color: Theme.of(context).cardTheme.color,
              borderRadius: BorderRadius.circular(18).copyWith(
                bottomLeft: const Radius.circular(4),
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withAlpha(10),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                _buildDot(0),
                const SizedBox(width: 6),
                _buildDot(1),
                const SizedBox(width: 6),
                _buildDot(2),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDot(int index) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        // Cada punto tiene un delay diferente para efecto de ola
        final delay = index * 0.2;
        final value = (_controller.value - delay) % 1.0;
        
        // Animación de scale: crece y decrece
        final scale = value < 0.5
            ? 1.0 + (value * 2) * 0.5  // Crece de 1.0 a 1.5
            : 1.5 - ((value - 0.5) * 2) * 0.5;  // Decrece de 1.5 a 1.0
        
        // Animación de opacity: más visible cuando está más grande
        final opacity = 0.3 + (scale - 1.0);
        
        return Transform.scale(
          scale: scale,
          child: Opacity(
            opacity: opacity.clamp(0.3, 1.0),
            child: Container(
              width: 8,
              height: 8,
              decoration: BoxDecoration(
                color: AppTheme.primaryBlue,
                shape: BoxShape.circle,
              ),
            ),
          ),
        );
      },
    );
  }
}

/// Versión compacta del typing indicator (sin avatar, para otras UIs)
class TypingIndicatorCompact extends StatefulWidget {
  final Color? dotColor;
  final double dotSize;

  const TypingIndicatorCompact({
    super.key,
    this.dotColor,
    this.dotSize = 8.0,
  });

  @override
  State<TypingIndicatorCompact> createState() => _TypingIndicatorCompactState();
}

class _TypingIndicatorCompactState extends State<TypingIndicatorCompact>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1400),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final color = widget.dotColor ?? AppTheme.primaryBlue;

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        _buildDot(0, color),
        const SizedBox(width: 6),
        _buildDot(1, color),
        const SizedBox(width: 6),
        _buildDot(2, color),
      ],
    );
  }

  Widget _buildDot(int index, Color color) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        final delay = index * 0.2;
        final value = (_controller.value - delay) % 1.0;
        
        final scale = value < 0.5
            ? 1.0 + (value * 2) * 0.5
            : 1.5 - ((value - 0.5) * 2) * 0.5;
        
        final opacity = 0.3 + (scale - 1.0);
        
        return Transform.scale(
          scale: scale,
          child: Opacity(
            opacity: opacity.clamp(0.3, 1.0),
            child: Container(
              width: widget.dotSize,
              height: widget.dotSize,
              decoration: BoxDecoration(
                color: color,
                shape: BoxShape.circle,
              ),
            ),
          ),
        );
      },
    );
  }
}
