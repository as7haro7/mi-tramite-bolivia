import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../providers/app_provider.dart';
import '../providers/assistant_provider.dart';
import '../widgets/source_card.dart';
import '../widgets/typing_indicator.dart';

class AssistantScreen extends StatefulWidget {
  final bool isBottomSheet;

  const AssistantScreen({super.key, this.isBottomSheet = false});

  @override
  State<AssistantScreen> createState() => _AssistantScreenState();
}

class _AssistantScreenState extends State<AssistantScreen> {
  final TextEditingController _textController = TextEditingController();
  final ScrollController _scrollController = ScrollController();

  final List<String> _quickPrompts = [
    '¿Qué necesito para sacar el NIT?',
    '¿Cómo renuevo mi Pasaporte?',
    '¿Costo de Licencia de Conducir?',
    '¿Registro SEPREC?',
  ];

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _handleSend(String text) {
    if (text.trim().isEmpty) return;
    final city = Provider.of<AppProvider>(context, listen: false).selectedCity;
    final provider = Provider.of<AssistantProvider>(context, listen: false);
    provider.sendMessage(text, city: city);
    _textController.clear();
    
    // Scroll inicial para mostrar mensaje del usuario
    _scrollToBottom();
    
    // Segundo scroll después de un frame para mostrar typing indicator
    Future.delayed(const Duration(milliseconds: 100), () {
      _scrollToBottom();
    });
  }

  Widget _buildFormattedText(String text, TextStyle baseStyle) {
    final List<InlineSpan> spans = [];
    final RegExp regExp = RegExp(r'\*\*(.*?)\*\*');
    int lastMatchEnd = 0;

    for (final Match match in regExp.allMatches(text)) {
      if (match.start > lastMatchEnd) {
        spans.add(TextSpan(text: text.substring(lastMatchEnd, match.start)));
      }
      final String boldText = match.group(1) ?? '';
      spans.add(TextSpan(
        text: boldText,
        style: baseStyle.copyWith(fontWeight: FontWeight.bold),
      ));
      lastMatchEnd = match.end;
    }

    if (lastMatchEnd < text.length) {
      spans.add(TextSpan(text: text.substring(lastMatchEnd)));
    }

    return SelectableText.rich(
      TextSpan(style: baseStyle, children: spans),
    );
  }

  @override
  void dispose() {
    _textController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final assistant = Provider.of<AssistantProvider>(context);

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        automaticallyImplyLeading: !widget.isBottomSheet,
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: const [
            Icon(Icons.auto_awesome_rounded, color: Colors.amber, size: 20),
            SizedBox(width: 8),
            Text('Asistente IA RAG', style: TextStyle(fontWeight: FontWeight.bold)),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh_rounded),
            tooltip: 'Nueva conversación',
            onPressed: () => assistant.clearChat(),
          ),
          if (widget.isBottomSheet)
            IconButton(
              icon: const Icon(Icons.close_rounded),
              onPressed: () => Navigator.pop(context),
            ),
        ],
      ),
      body: Column(
        children: [
          // iOS AI Disclaimer Banner
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: BoxDecoration(
              color: Colors.amber.withAlpha(25),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.amber.withAlpha(60)),
            ),
            child: Row(
              children: const [
                Icon(Icons.shield_outlined, size: 16, color: AppTheme.accentAmber),
                SizedBox(width: 8),
                Expanded(
                  child: Text(
                    'La IA procesa datos oficiales. Verifica siempre las fuentes adjuntas.',
                    style: TextStyle(fontSize: 11.5, fontWeight: FontWeight.w600),
                  ),
                ),
              ],
            ),
          ),

          // Chat Messages List
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: const EdgeInsets.all(16),
              itemCount: assistant.messages.length + (assistant.isLoading ? 1 : 0),
              itemBuilder: (context, index) {
                // Mostrar typing indicator al final cuando está cargando
                if (index == assistant.messages.length) {
                  return const TypingIndicator();
                }
                
                final msg = assistant.messages[index];
                return Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: Row(
                    mainAxisAlignment:
                        msg.isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      if (!msg.isUser) ...[
                        Container(
                          width: 32,
                          height: 32,
                          decoration: const BoxDecoration(
                            gradient: AppTheme.primaryGradient,
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(Icons.auto_awesome_rounded, color: Colors.white, size: 16),
                        ),
                        const SizedBox(width: 8),
                      ],
                      Flexible(
                        child: Container(
                          padding: const EdgeInsets.all(14),
                          decoration: BoxDecoration(
                            color: msg.isUser
                                ? AppTheme.iosBlue
                                : Theme.of(context).cardTheme.color,
                            borderRadius: BorderRadius.circular(18).copyWith(
                              bottomRight: msg.isUser ? const Radius.circular(4) : null,
                              bottomLeft: !msg.isUser ? const Radius.circular(4) : null,
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withAlpha(10),
                                blurRadius: 4,
                                offset: const Offset(0, 2),
                              ),
                            ],
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              _buildFormattedText(
                                msg.content,
                                TextStyle(
                                  color: msg.isUser
                                      ? Colors.white
                                      : Theme.of(context).colorScheme.onSurface,
                                  fontSize: 14,
                                  height: 1.4,
                                ),
                              ),
                              if (!msg.isUser && msg.fuentes.isNotEmpty) ...[
                                const SizedBox(height: 10),
                                const Text(
                                  'Fuentes Oficiales Citadas:',
                                  style: TextStyle(
                                      fontWeight: FontWeight.bold,
                                      fontSize: 11.5,
                                      color: AppTheme.iosTeal),
                                ),
                                ...msg.fuentes.map((f) => SourceCard(title: f.titulo, url: f.url)),
                              ],
                              if (!msg.isUser && msg.id != 'welcome') ...[
                                const SizedBox(height: 6),
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.end,
                                  children: [
                                    const Text('¿Fue útil?', style: TextStyle(fontSize: 11, color: Colors.grey)),
                                    IconButton(
                                      constraints: const BoxConstraints(),
                                      padding: const EdgeInsets.symmetric(horizontal: 4),
                                      icon: Icon(
                                        Icons.thumb_up_rounded,
                                        size: 15,
                                        color: msg.feedback == true ? AppTheme.successGreen : Colors.grey,
                                      ),
                                      onPressed: () => assistant.sendFeedback(msg.id, true),
                                    ),
                                    IconButton(
                                      constraints: const BoxConstraints(),
                                      padding: const EdgeInsets.symmetric(horizontal: 4),
                                      icon: Icon(
                                        Icons.thumb_down_rounded,
                                        size: 15,
                                        color: msg.feedback == false ? AppTheme.alertRed : Colors.grey,
                                      ),
                                      onPressed: () => assistant.sendFeedback(msg.id, false),
                                    ),
                                  ],
                                ),
                              ],
                            ],
                          ),
                        ),
                      ),
                      if (msg.isUser) ...[
                        const SizedBox(width: 8),
                        const CircleAvatar(
                          radius: 16,
                          backgroundColor: AppTheme.iosBlue,
                          child: Icon(Icons.person_rounded, color: Colors.white, size: 18),
                        ),
                      ],
                    ],
                  ),
                );
              },
            ),
          ),

          if (assistant.isLoading)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 6),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  SizedBox(width: 14, height: 14, child: CircularProgressIndicator(strokeWidth: 2)),
                  SizedBox(width: 8),
                  Text('Consultando información verificada...', style: TextStyle(fontSize: 12, color: Colors.grey)),
                ],
              ),
            ),

          // Quick Prompts Chips
          SizedBox(
            height: 36,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 12),
              itemCount: _quickPrompts.length,
              itemBuilder: (context, index) {
                final prompt = _quickPrompts[index];
                return Padding(
                  padding: const EdgeInsets.only(right: 6),
                  child: ActionChip(
                    label: Text(prompt, style: const TextStyle(fontSize: 12)),
                    onPressed: () => _handleSend(prompt),
                  ),
                );
              },
            ),
          ),
          const SizedBox(height: 8),

          // Input Bar (iOS HIG feel)
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: BoxDecoration(
              color: Theme.of(context).cardTheme.color,
              border: const Border(top: BorderSide(color: AppTheme.iosLightBorder, width: 0.8)),
            ),
            child: Row(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14),
                    decoration: BoxDecoration(
                      color: Theme.of(context).scaffoldBackgroundColor,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: TextField(
                      controller: _textController,
                      onSubmitted: _handleSend,
                      decoration: const InputDecoration(
                        hintText: 'Consulta sobre un trámite...',
                        border: InputBorder.none,
                        isDense: true,
                        contentPadding: EdgeInsets.symmetric(vertical: 10),
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton(
                  icon: const Icon(Icons.arrow_upward_rounded, color: Colors.white),
                  style: IconButton.styleFrom(
                    backgroundColor: AppTheme.iosBlue,
                    padding: const EdgeInsets.all(8),
                    minimumSize: Size.zero,
                  ),
                  onPressed: () => _handleSend(_textController.text),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
